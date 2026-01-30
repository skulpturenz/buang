package fitnesstests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"
	testutils "skulpture/buang/tests/utils"
	"skulpture/buang/utils/compensations"
	"strconv"
	"sync"
	"testing"

	"github.com/go-git/go-billy/v6/osfs"
	httpbackend "github.com/go-git/go-git/v6/backend/http"
	"github.com/go-git/go-git/v6/plumbing/transport"
	httpplumbing "github.com/go-git/go-git/v6/plumbing/transport/http"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/stretchr/testify/require"
)

func TestDeploymentsFitness(t *testing.T) {
	ctx := t.Context()

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(ctx)
	require.NoError(t, err)
	defer temporalPgCleanup(ctx)

	MIN_DEPLOYMENTS := 60 // at once, no build, nginx, just deploy

	// a project can only be linked to one repository
	// project to git repo is 1-1
	// so setup a git server and serve "multiple" repositories
	gitServerUrl, gitCleanup := setupLocalGitServer(t, MIN_DEPLOYMENTS)
	defer gitCleanup(ctx)

	compensations := compensations.New()
	defer compensations.Compensate(ctx)

	testApp, cleanup := testutils.Setup(ctx, *temporalPgConfig)
	compensations.AddCompensation(cleanup)

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)
	var wg sync.WaitGroup

	createProject := func(repoUrl string) int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:    repoUrl,
			RequiresAuthn: false,
			ComposePath:   "compose.yaml",
		}

		body, err := json.Marshal(createProjectReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, postProjectsUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		projectId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return projectId
	}

	createDeployment := func(projectId int64) int64 {
		postDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment?waitForDeployment=true", baseUrl, projectId)
		createDeploymentReq := deploymentshandlers.CreateDeploymentRequest{
			Branch:            "master",
			Sha:               "e6792e4fe8a66de90b0945fa9d38f0b25149bd00",
			ServiceEntrypoint: "web:80",
			Env: map[string]any{
				"HELLO": "WORLD",
			},
		}

		body, err := json.Marshal(createDeploymentReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, postDeploymentUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		deploymentId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return deploymentId
	}

	deleteProject := func(projectId int64) { // instead of buang deployments because this blocks
		deleteProjectUrl := fmt.Sprintf("%v/project/%v", baseUrl, projectId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteProjectUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	}

	for i := range MIN_DEPLOYMENTS {
		wg.Add(1)

		deploy := func(idx int) {
			defer wg.Done()
			repoUrl := fmt.Sprintf("%s/repo-%d.git", gitServerUrl, idx)
			projectId := createProject(repoUrl)
			createDeployment(projectId)
			compensations.AddCompensation(func(ctx context.Context) {
				deleteProject(projectId) // blocks
			})
		}

		go deploy(i)
	}

	wg.Wait()
}

func setupLocalGitServer(t *testing.T, x int) (string, func(context.Context)) {
	ctx := t.Context()

	compensations := compensations.New()

	tmp, err := os.MkdirTemp("/tmp", "buang-fitness-test-*")
	require.NoError(t, err)
	compensations.AddCompensation(func(ctx context.Context) {
		os.RemoveAll(tmp)
	})

	baseRepoPath := filepath.Join(tmp, "base")
	_, err = git.PlainCloneContext(ctx, baseRepoPath, &git.CloneOptions{
		URL:             "https://github.com/skulpturenz/buangtest", // TODO: from repo vars,
		NoCheckout:      false,
		InsecureSkipTLS: true,
		ReferenceName:   plumbing.NewBranchReferenceName("master"),
		SingleBranch:    true,
		Auth: &httpplumbing.BasicAuth{
			Username: os.Getenv("BUANG_TEST_GITHUB_USER"), // TODO: from repo secrets
			Password: os.Getenv("BUANG_TEST_GITHUB_PAT"),  // TODO: from repo secrets
		},
		Bare: true,
	})
	require.NoError(t, err)

	for i := range x {
		// reference:
		// - https://pkg.go.dev/github.com/go-git/go-git/v6/backend/http#NewBackend
		// - https://pkg.go.dev/github.com/go-git/go-billy/v5/osfs#New
		repo := fmt.Sprintf("repo-%d.git", i)
		repoPath := filepath.Join(tmp, repo)

		cmd := exec.Command("cp", "-r", baseRepoPath, repoPath)
		err := cmd.Run()
		require.NoError(t, err)

		cmd = exec.Command("git", "--git-dir", repoPath, "update-server-info")
		err = cmd.Run()
		require.NoError(t, err)
	}

	server := httptest.NewServer(httpbackend.NewBackend(transport.NewFilesystemLoader(osfs.New(tmp), false)))
	compensations.AddCompensation(func(ctx context.Context) {
		server.Close()
	})

	return server.URL, compensations.Compensate
}
