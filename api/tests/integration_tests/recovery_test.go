package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"
	proxiesrecoverycomposedown "skulpture/buang/tests/proxies/recovery_compose_down_test"
	proxiesrecoverygitclone "skulpture/buang/tests/proxies/recovery_git_clone_test"
	proxiesrecoveryimagepull "skulpture/buang/tests/proxies/recovery_image_pull_test"
	testutils "skulpture/buang/tests/utils"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestImagePullError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(ctx)
	require.NoError(t, err)
	defer dbosCleanup(ctx)

	temporalSqliteConfig, temporalSqliteCleanup, err := testutils.CreateSqliteTemporal(ctx)
	require.NoError(t, err)
	defer temporalSqliteCleanup(ctx)

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(ctx)
	require.NoError(t, err)
	defer temporalPgCleanup(ctx)

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalSqliteConfig,
		"TemporalPg":     *temporalPgConfig,
	}

	retry := testutils.NewRetry(3, 500*time.Millisecond)

	for k, v := range scenarios {
		retry.Retry(t, func(t *testing.T) {
			t.Run(k, func(t *testing.T) {
				imagePullError(t, v)
			})
		})
	}
}

func TestComposeDownError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(ctx)
	require.NoError(t, err)
	defer dbosCleanup(ctx)

	temporalSqliteConfig, temporalSqliteCleanup, err := testutils.CreateSqliteTemporal(ctx)
	require.NoError(t, err)
	defer temporalSqliteCleanup(ctx)

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(ctx)
	require.NoError(t, err)
	defer temporalPgCleanup(ctx)

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalSqliteConfig,
		"TemporalPg":     *temporalPgConfig,
	}

	retry := testutils.NewRetry(3, 500*time.Millisecond)

	for k, v := range scenarios {
		retry.Retry(t, func(t *testing.T) {
			t.Run(k, func(t *testing.T) {
				composeDownError(t, v)
			})
		})
	}
}

func TestGitCloneError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(ctx)
	require.NoError(t, err)
	defer dbosCleanup(ctx)

	temporalSqliteConfig, temporalSqliteCleanup, err := testutils.CreateSqliteTemporal(ctx)
	require.NoError(t, err)
	defer temporalSqliteCleanup(ctx)

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(ctx)
	require.NoError(t, err)
	defer temporalPgCleanup(ctx)

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalSqliteConfig,
		"TemporalPg":     *temporalPgConfig,
	}

	retry := testutils.NewRetry(3, 500*time.Millisecond)

	for k, v := range scenarios {
		retry.Retry(t, func(t *testing.T) {
			t.Run(k, func(t *testing.T) {
				gitCloneError(t, v)
			})
		})
	}
}

func imagePullError(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

	testApp, cleanup := testutils.Setup(ctx, config, proxiesrecoveryimagepull.New())
	defer cleanup(ctx)

	githubPat := os.Getenv("BUANG_TEST_GITHUB_PAT")
	username := os.Getenv("BUANG_TEST_GITHUB_USER")

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)

	createProject := func() int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:    "https://github.com/skulpturenz/buangtest", // TODO: from repo vars
			RequiresAuthn: true,
			Username:      &username,      // TODO: from repo secrets
			Password:      &githubPat,     // TODO: from repo secrets
			ComposePath:   "compose.yaml", // TODO: from repo vars
		}

		body, err := json.Marshal(createProjectReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, postProjectsUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create project")
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		projectId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return projectId
	}

	createDeployment := func(projectId int64) int64 {
		postDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment?waitForDeployment=true", baseUrl, projectId)
		createDeploymentReq := deploymentshandlers.CreateDeploymentRequest{
			Branch:            "master",                                   // TODO: from repo vars
			Sha:               "4b591d8f2b1b2a28645d98911d84b2eb68f7142d", // TODO: from repo vars
			ServiceEntrypoint: "web:80",                                   // TODO: from repo vars
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

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create deployment")
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		deploymentId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return deploymentId
	}

	buangDeployment := func(projectId int64, deploymentId int64) {
		buangDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment/%v", baseUrl, projectId, deploymentId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, buangDeploymentUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusNoContent, res.StatusCode, "failed to buang deployment")

		time.Sleep(500 * time.Millisecond)
	}

	projectId := createProject()
	deploymentId := createDeployment(projectId)
	buangDeployment(projectId, deploymentId)
}

func composeDownError(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

	testApp, cleanup := testutils.Setup(ctx, config, proxiesrecoverycomposedown.New())
	defer cleanup(ctx)

	githubPat := os.Getenv("BUANG_TEST_GITHUB_PAT")
	username := os.Getenv("BUANG_TEST_GITHUB_USER")

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)

	createProject := func() int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:    "https://github.com/skulpturenz/buangtest", // TODO: from repo vars
			RequiresAuthn: true,
			Username:      &username,      // TODO: from repo secrets
			Password:      &githubPat,     // TODO: from repo secrets
			ComposePath:   "compose.yaml", // TODO: from repo vars
		}

		body, err := json.Marshal(createProjectReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, postProjectsUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create project")
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		projectId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return projectId
	}

	createDeployment := func(projectId int64) int64 {
		postDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment?waitForDeployment=true", baseUrl, projectId)
		createDeploymentReq := deploymentshandlers.CreateDeploymentRequest{
			Branch:            "master",                                   // TODO: from repo vars
			Sha:               "4b591d8f2b1b2a28645d98911d84b2eb68f7142d", // TODO: from repo vars
			ServiceEntrypoint: "web:80",                                   // TODO: from repo vars
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

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create deployment")
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

	projectId := createProject()
	createDeployment(projectId)
	deleteProject(projectId)
}

func gitCloneError(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

	testApp, cleanup := testutils.Setup(ctx, config, proxiesrecoverygitclone.New())
	defer cleanup(ctx)

	githubPat := os.Getenv("BUANG_TEST_GITHUB_PAT")
	username := os.Getenv("BUANG_TEST_GITHUB_USER")

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)

	createProject := func() int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:    "https://github.com/skulpturenz/buangtest", // TODO: from repo vars
			RequiresAuthn: true,
			Username:      &username,      // TODO: from repo secrets
			Password:      &githubPat,     // TODO: from repo secrets
			ComposePath:   "compose.yaml", // TODO: from repo vars
		}

		body, err := json.Marshal(createProjectReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, postProjectsUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create project")
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		projectId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return projectId
	}

	createDeployment := func(projectId int64) int64 {
		postDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment?waitForDeployment=true", baseUrl, projectId)
		createDeploymentReq := deploymentshandlers.CreateDeploymentRequest{
			Branch:            "master",                                   // TODO: from repo vars
			Sha:               "4b591d8f2b1b2a28645d98911d84b2eb68f7142d", // TODO: from repo vars
			ServiceEntrypoint: "web:80",                                   // TODO: from repo vars
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

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to create deployment")
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

	projectId := createProject()
	createDeployment(projectId)
	deleteProject(projectId)
}
