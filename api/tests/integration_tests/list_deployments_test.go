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
	testutils "skulpture/buang/tests/utils"
	"skulpture/buang/utils/compensations"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/negrel/assert"
	"github.com/stretchr/testify/require"
)

func TestListDeployments(t *testing.T) {
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
				listDeployments(t, v)
			})
		})
	}
}

func TestSearchParams(t *testing.T) {
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
				searchParams(t, v)
			})
		})
	}
}

func TestDeletedProjects(t *testing.T) {
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
				deletedProject(t, v)
			})
		})
	}
}

func listDeployments(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

	compensations := compensations.New()
	defer compensations.Compensate(ctx)

	testApp, cleanup := testutils.Setup(ctx, config)
	compensations.AddCompensation(func(ctx context.Context) {
		cleanup(ctx)
	})

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

	buangDeploymenbt := func(projectId int64, deploymentId int64) {
		buangDeploymentUrl := fmt.Sprintf("%v/project/%v/deployment/%v", baseUrl, projectId, deploymentId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, buangDeploymentUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode, "failed to buang deployment")
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

	deployments := listDeployments(projectId)

	require.Len(t, deployments, 0)

	firstDeploymentId := createDeployment(projectId)
	compensations.AddCompensation(func(ctx context.Context) {
		buangDeploymenbt(projectId, firstDeploymentId)
		time.Sleep(500 * time.Millisecond) // async workflow
	})

	secondDeploymentId := createDeployment(projectId)
	compensations.AddCompensation(func(ctx context.Context) {
		buangDeploymenbt(projectId, secondDeploymentId)
		time.Sleep(500 * time.Millisecond) // async workflow
	})

	deployments = listDeployments(projectId)

	require.Len(t, deployments, 2)
	require.Equal(t, projectId, deployments[0].ProjectID)

	unsorted := []int64{} // assume desc
	for _, i := range deployments {
		unsorted = append(unsorted, i.ID)
	}

	sorted := unsorted // asc
	slices.Sort(sorted)

	slices.Reverse(unsorted)
	require.Equal(t, unsorted, sorted)
}

func searchParams(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

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

	listDeployments := func(projectId int64, limit int, page int) []projectshandlers.ListAllDeploymentsItem {
		listDeploymentsUrl := fmt.Sprintf("%v/project/%v/deployments?limit=%v&page=%v", baseUrl, projectId, limit, page)
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

	deleteProject := func(projectId int64) {
		deleteProjectUrl := fmt.Sprintf("%v/project/%v", baseUrl, projectId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteProjectUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	}

	projectId := createProject()
	firstDeployment := createDeployment(projectId)
	secondDeployment := createDeployment(projectId)
	thirdDeployment := createDeployment(projectId)

	limitedDeployments := listDeployments(projectId, 2, 1)
	require.Len(t, limitedDeployments, 2)
	require.Equal(t, thirdDeployment, limitedDeployments[0].ID)
	require.Equal(t, secondDeployment, limitedDeployments[1].ID)

	pagedDeployments := listDeployments(projectId, 2, 2)
	require.Len(t, pagedDeployments, 1)
	require.Equal(t, firstDeployment, pagedDeployments[0].ID)

	deleteProject(projectId)
}

func deletedProject(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx := t.Context()

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

	deleteProject := func(projectId int64) {
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

	deployments := listDeployments(projectId)
	require.Len(t, deployments, 0)
}
