package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	projectshandlers "skulpture/buang/handlers/projects"
	testutils "skulpture/buang/integration_tests/utils"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListProjectsDescendingOrder(t *testing.T) {
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
				listProjectsDescendingOrder(t, v)
			})
		})
	}
}

func TestListProjectsExcludesDeleted(t *testing.T) {
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
				listProjectsExcludesDeleted(t, v)
			})
		})
	}
}

func listProjectsDescendingOrder(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	testApp, cleanup := testutils.Setup(ctx, config)
	defer cleanup(ctx)

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)

	createProject := func(repository string) int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:  repository,
			ComposePath: "compose.yaml",
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

	_ = createProject("https://github.com/repo1")
	_ = createProject("https://github.com/repo2")
	_ = createProject("https://github.com/repo3")

	getProjectsUrl := fmt.Sprintf("%v/projects", baseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getProjectsUrl, nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var projects []projectshandlers.ListAllProjectsItem
	err = json.NewDecoder(res.Body).Decode(&projects)
	require.NoError(t, err)

	unsorted := []int64{} // assume desc
	for _, i := range projects {
		unsorted = append(unsorted, i.Id)
	}

	sorted := unsorted // asc
	slices.Sort(sorted)

	slices.Reverse(unsorted)
	require.Equal(t, unsorted, sorted)
}

func listProjectsExcludesDeleted(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Minute)
	defer cancel()

	testApp, cleanup := testutils.Setup(ctx, config)
	defer cleanup(ctx)

	baseUrl := fmt.Sprintf("%v/api/v1", *testApp.Url)

	createProject := func(repository string) int64 {
		postProjectsUrl := fmt.Sprintf("%v/project", baseUrl)
		createProjectReq := projectshandlers.CreateProjectRequest{
			Repository:  repository,
			ComposePath: "compose.yaml",
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

	deleteProject := func(projectId int64) {
		deleteProjectUrl := fmt.Sprintf("%v/project/%v", baseUrl, projectId)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteProjectUrl, nil)
		require.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusNoContent, res.StatusCode, "failed to delete project")
	}

	p1 := createProject("https://github.com/repo1")
	p2 := createProject("https://github.com/repo2")
	p3 := createProject("https://github.com/repo3")

	deleteProject(p2)

	getProjectsUrl := fmt.Sprintf("%v/projects", baseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getProjectsUrl, nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var projects []projectshandlers.ListAllProjectsItem
	err = json.NewDecoder(res.Body).Decode(&projects)
	require.NoError(t, err)

	require.Len(t, projects, 2)

	ids := []int64{}
	for _, p := range projects {
		ids = append(ids, p.Id)
	}

	require.Contains(t, ids, p1)
	require.Contains(t, ids, p3)
	require.NotContains(t, ids, p2)
}
