package fitnesstests

import (
	"archive/tar"
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
	deploymentscomponent "skulpture/buang/components/deployments"
	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"
	testutils "skulpture/buang/tests/utils"
	"skulpture/buang/utils/compensations"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
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

	docker := testApp.GetHttpApplication().Services.Docker
	createTraefik := func() (string, string) {
		reader, err := docker.ImagePull(ctx, "traefik", image.PullOptions{})
		require.NoError(t, err)
		io.Copy(io.Discard, reader)

		web, err := nat.NewPort("tcp", "80")
		require.NoError(t, err)

		traefikConfig := container.Config{
			Image: "traefik",
			Cmd: []string{
				"--providers.file.directory=/app/deployments",
				"--providers.file.watch=true",
				"--entryPoints.web.address=:80",
			},
			ExposedPorts: nat.PortSet{
				web: struct{}{},
			},
		}
		traefikHostConfig := container.HostConfig{
			PortBindings: nat.PortMap{
				web: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "0", // random
					},
				},
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("traefik-vol-%v", time.Now().Nanosecond()),
					Target: "/app/deployments",
				},
			},
		}

		containerName := fmt.Sprintf("traefik_%v", time.Now().Nanosecond())
		traefik, err := docker.ContainerCreate(ctx, &traefikConfig, &traefikHostConfig, nil, nil, containerName)
		require.NoError(t, err)

		err = docker.ContainerStart(ctx, traefik.ID, container.StartOptions{})
		require.NoError(t, err)
		compensations.AddCompensation(func(ctx context.Context) {
			docker.ContainerRemove(ctx, traefik.ID, container.RemoveOptions{Force: true})
		})

		inspect, err := docker.ContainerInspect(ctx, traefik.ID)
		require.NoError(t, err)
		ports := inspect.NetworkSettings.Ports[web]
		require.NotEmpty(t, ports)

		return traefik.ID, ports[0].HostPort
	}

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

	assertProxy := func(projectId int64, deploymentId int64, containerId string, port string) {
		s := testApp.GetHttpApplication().Services.ToAppApplicationServices()

		findDeploymentParams := deploymentscomponent.FindDeploymentByIdParams{
			ID:        deploymentId,
			ProjectId: projectId,
		}

		deployment, err := findDeploymentParams.Exec(ctx, &s)
		require.NoError(t, err)

		deploymentConfigPath := deploymentscomponent.GetDeploymentPath(deploymentscomponent.GetDeploymentPathParams{
			ProjectId:    projectId,
			DeploymentId: deploymentId,
			Branch:       deployment.Deployment.GetBranch(),
			Sha:          deployment.Deployment.GetSha(),
			DeployedAt:   *deployment.Deployment.GetDeployedAt(),
		})

		configContent, err := os.ReadFile(deploymentConfigPath)
		require.NoError(t, err)

		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		header := &tar.Header{
			Name: filepath.Base(deploymentConfigPath),
			Mode: 0644, // creator: rw, others: r
			Size: int64(len(configContent)),
		}
		err = tw.WriteHeader(header)
		require.NoError(t, err)
		_, err = tw.Write(configContent)
		require.NoError(t, err)
		err = tw.Close()
		require.NoError(t, err)

		err = docker.CopyToContainer(ctx, containerId, "/app/deployments", &buf, container.CopyToContainerOptions{})
		require.NoError(t, err)

		projectName := deploymentscomponent.GetProjectName(deploymentscomponent.GetProjectNameParams{
			ProjectId:    projectId,
			DeploymentId: deploymentId,
			Branch:       deployment.Deployment.GetBranch(),
			Sha:          deployment.Deployment.GetSha(),
			DeployedAt:   *deployment.Deployment.GetDeployedAt(),
		})
		filters := filters.NewArgs()
		filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
		testServiceContainers, err := docker.ContainerList(ctx, container.ListOptions{
			Filters: filters,
		})
		require.NoError(t, err)
		require.Len(t, testServiceContainers, 1)

		var networkID string
		for _, n := range testServiceContainers[0].NetworkSettings.Networks {
			if n.NetworkID == "bridge" {
				continue
			}

			networkID = n.NetworkID
		}
		require.NotEmpty(t, networkID)

		err = docker.NetworkConnect(ctx, networkID, containerId, &network.EndpointSettings{})
		if err != nil {
			t.Logf("assert proxy err: %v", err)
		}

		time.Sleep(3 * time.Second)

		instanceUrl := fmt.Sprintf("http://localhost:%v%v", port, *deployment.Deployment.GetUrl())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, instanceUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)

		logs, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.Contains(t, string(logs), "nginx")
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

	traefikId, traefikPort := createTraefik()
	for i := range MIN_DEPLOYMENTS {
		wg.Add(1)

		deploy := func(idx int) {
			defer wg.Done()
			repoUrl := fmt.Sprintf("%s/repo-%d.git", gitServerUrl, idx)
			projectId := createProject(repoUrl)
			deploymentId := createDeployment(projectId)
			assertProxy(projectId, deploymentId, traefikId, traefikPort)
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
