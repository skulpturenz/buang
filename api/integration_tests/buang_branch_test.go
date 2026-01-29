package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"
	testutils "skulpture/buang/integration_tests/utils"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuangBranch(t *testing.T) {
	t.Parallel()

	dbosConfig, dbosCleanup, err := testutils.CreateDbos(t.Context())
	require.NoError(t, err)
	defer dbosCleanup(t.Context())

	temporalSqliteConfig, temporalSqliteCleanup, err := testutils.CreateSqliteTemporal(t.Context())
	require.NoError(t, err)
	defer temporalSqliteCleanup(t.Context())

	temporalPgConfig, temporalPgCleanup, err := testutils.CreatePgTemporal(t.Context())
	require.NoError(t, err)
	defer temporalPgCleanup(t.Context())

	scenarios := map[string]testutils.DurableExecutorConfiguration{
		"DBOS":           *dbosConfig,
		"TemporalSqlite": *temporalSqliteConfig,
		"TemporalPg":     *temporalPgConfig,
	}

	retry := testutils.NewRetry(3, 500*time.Millisecond)

	for k, v := range scenarios {
		retry.Retry(t, func(t *testing.T) {
			t.Run(k, func(t *testing.T) {
				buangBranch(t, v)
			})
		})
	}
}

func buangBranch(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	testApp, cleanup := testutils.Setup(ctx, config)
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
			Branch:            "master",                                   // TODO: from repo vars
			Sha:               "e6792e4fe8a66de90b0945fa9d38f0b25149bd00", // TODO: from repo vars
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

		require.Equal(t, http.StatusOK, res.StatusCode)
		respBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		deploymentId, err := strconv.ParseInt(string(respBody), 10, 64)
		require.NoError(t, err)

		return deploymentId
	}

	buangBranch := func(projectId int64) {
		buangBranchUrl := fmt.Sprintf("%v/project/%v/branch", baseUrl, projectId)
		buangBranchReq := projectshandlers.BuangBranchRequest{
			Branch: "master",
		}

		body, err := json.Marshal(buangBranchReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, buangBranchUrl, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	}

	listDeployments := func(projectId int64) []projectshandlers.ListAllDeploymentsItem {
		listDeploymentsUrl := fmt.Sprintf("%v/project/%v/deployments", baseUrl, projectId)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, listDeploymentsUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		var deployments []projectshandlers.ListAllDeploymentsItem
		err = json.NewDecoder(res.Body).Decode(&deployments)
		require.NoError(t, err)

		return deployments
	}

	projectId := createProject()
	deploymentId := createDeployment(projectId)
	buangBranch(projectId)

	time.Sleep(5 * time.Second) // async workflow so returns immediately

	deployments := listDeployments(projectId)
	idx := slices.IndexFunc(deployments, func(deployment projectshandlers.ListAllDeploymentsItem) bool {
		return deployment.ID == deploymentId
	})
	require.Equal(t, deployments[idx].Status, int16(enumsdeploymentstatus.Buang))
}
