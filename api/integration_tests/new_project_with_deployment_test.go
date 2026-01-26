package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"skulpture/buang/db"
	testutils "skulpture/buang/integration_tests/utils"
	"strconv"
	"testing"
	"time"

	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"

	"github.com/docker/docker/client"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProjectWithDeployment(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(ctx)
	require.NoError(t, err)
	defer dbosCleanup(ctx)

	temporalConfig, temporalCleanup, err := testutils.CreateSqliteTemporal(ctx)
	require.NoError(t, err)
	defer temporalCleanup(ctx)

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalConfig,
	}

	for k, v := range scenarios {
		t.Run(k, func(t *testing.T) {
			createNewProjectWithDeployment(t, v)
		})
	}
}

func createNewProjectWithDeployment(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	dbCfg := db.DbConfig{
		Type:             config.DbType,
		ConnectionString: config.DbConnectionString,
	}
	queries, dbCleanup, err := dbCfg.New(ctx)
	require.NoError(t, err)
	defer dbCleanup(ctx)

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)
	defer docker.Close()

	decoder := schema.NewDecoder()
	encoder := schema.NewEncoder()

	appServices := testutils.TestApplicationServices{
		Queries:              &queries,
		GorillaSchemaDecoder: decoder,
		GorillaSchemaEncoder: encoder,
		Docker:               docker,
	}

	appConfig := testutils.TestApplicationConfig{
		Services:              appServices,
		DurableExecutorConfig: config,
	}

	testApp, appCleanup := testutils.CreateApp(ctx, appConfig)
	defer appCleanup(ctx)

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
		postDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment", baseUrl, projectId)
		createDeploymentReq := deploymentshandlers.CreateDeploymentRequest{
			Branch:            "master",                                   // TODO: from repo vars
			Sha:               "e6792e4fe8a66de90b0945fa9d38f0b25149bd00", // TODO: from repo vars
			ServiceEntrypoint: "web:80",                                   // TODO: from repo vars
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

	getDeploymentLogs := func(projectId int64, deploymentId int64) {
		getDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment/%v", baseUrl, projectId, deploymentId)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getDeploymentUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		var deployment deploymentshandlers.FindDeploymentByIdResponse
		err = json.NewDecoder(res.Body).Decode(&deployment)
		require.NoError(t, err)

		getDeploymentLogsUrl := fmt.Sprintf("%v/project/%v/deployment/%v/logs", baseUrl, projectId, deploymentId)
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, getDeploymentLogsUrl, nil)
		require.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode, "failed to get deployment logs")
		logs, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.Contains(t, string(logs), fmt.Sprintf("Checked out commit %v, branch %v", deployment.Sha, deployment.Branch))
		require.Contains(t, string(logs), "nginx")
	}

	// TODO: check new container

	buangDeploymenbt := func(projectId int64, deploymentId int64) {
		buangDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment/%v", baseUrl, projectId, deploymentId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, buangDeploymentUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode, "failed to buang deployment")
		time.Sleep(500 * time.Millisecond) // allow some time to cleanup
	}

	projectId := createProject()
	deploymentId := createDeployment(projectId)
	getDeploymentLogs(projectId, deploymentId)
	buangDeploymenbt(projectId, deploymentId)
}
