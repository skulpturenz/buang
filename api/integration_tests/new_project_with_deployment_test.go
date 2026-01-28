package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	deploymentscomponent "skulpture/buang/components/deployments"
	projectscomponent "skulpture/buang/components/projects"
	constantsenvs "skulpture/buang/constants/envs"
	testutils "skulpture/buang/integration_tests/utils"
	"strconv"
	"testing"
	"time"

	deploymentshandlers "skulpture/buang/handlers/deployments"
	projectshandlers "skulpture/buang/handlers/projects"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/negrel/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProjectWithDeployment(t *testing.T) {
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
				createNewProjectWithDeployment(t, v)
			})
		})
	}
}

func createNewProjectWithDeployment(t *testing.T, config testutils.DurableExecutorConfiguration) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	testApp, cleanup := testutils.Setup(ctx, config)
	defer cleanup(ctx)

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

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

	assertDeployment := func(projectId int64, deploymentId int64) {
		const EXPECTED_SERVICES = 1
		var TRAEFIK_DYNAMIC_CONFIG = constantsenvs.BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR.Value()

		s := testApp.GetHttpApplication().Services.ToAppApplicationServices()

		findDeploymentParams := deploymentscomponent.FindDeploymentByIdParams{
			ID:        deploymentId,
			ProjectId: projectId,
		}

		deployment, err := findDeploymentParams.Exec(ctx, &s)
		require.NoError(t, err)

		require.NotNil(t, deployment.Deployment.GetDeployedAt())

		findProjectParams := projectscomponent.FindProjectByIdParams{
			ID: deploymentId,
		}

		p, err := findProjectParams.Exec(ctx, &s)
		require.NoError(t, err)

		require.DirExists(t, TRAEFIK_DYNAMIC_CONFIG)

		sha := fmt.Sprintf("%.*s", 8, deployment.Deployment.GetSha())
		projectName := fmt.Sprintf("%v_%v_%v_%v_%v",
			p.Project.GetId(),
			deployment.Deployment.GetId(),
			deployment.Deployment.GetBranch(),
			sha,
			deployment.Deployment.GetDeployedAt().UnixMilli())
		deploymentConfigPath := filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("buang-%v.yaml", projectName))
		require.FileExists(t, deploymentConfigPath)

		cli, err := command.NewDockerCli(command.WithAPIClient(docker))
		require.NoError(t, err)
		err = cli.Initialize(&flags.ClientOptions{})
		require.NoError(t, err)

		svc, err := compose.NewComposeService(cli, compose.WithPrompt(compose.AlwaysOkPrompt()))
		require.NoError(t, err)

		composeProject, err := svc.LoadProject(ctx, api.ProjectLoadOptions{
			ConfigPaths: []string{filepath.Join(*deployment.Deployment.GetClonePath(), p.Project.GetComposePath())},
			ProjectName: projectName,
		})
		require.NoError(t, err)

		require.Len(t, composeProject.AllServices(), EXPECTED_SERVICES)

		filters := filters.NewArgs()
		filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
		containers, err := docker.ContainerList(ctx, container.ListOptions{
			Filters: filters,
		})
		require.NoError(t, err)

		for _, v := range containers {
			i, err := docker.ContainerInspect(ctx, v.ID)
			require.NoError(t, err)

			envs := types.NewMappingWithEquals(i.Config.Env)

			require.Equal(t, *envs["HELLO"], "WORLD")
			require.Equal(t, *envs["BUANG_DEPLOYMENT_PATH"], *deployment.Deployment.GetUrl())
		}
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

	projectId := createProject()
	deploymentId := createDeployment(projectId)
	getDeploymentLogs(projectId, deploymentId)
	assertDeployment(projectId, deploymentId)
	buangDeploymenbt(projectId, deploymentId)
}
