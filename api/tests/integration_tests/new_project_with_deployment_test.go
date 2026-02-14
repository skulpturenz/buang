package integrationtests

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	deploymentscomponent "skulpture/buang/components/deployments"
	projectscomponent "skulpture/buang/components/projects"
	constantsenvs "skulpture/buang/constants/envs"
	testutils "skulpture/buang/tests/utils"
	"skulpture/buang/utils/compensations"
	"slices"
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
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
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
	ctx := t.Context()

	compensations := compensations.New()
	defer compensations.Compensate(ctx)

	testApp, cleanup := testutils.Setup(ctx, config, nil)
	compensations.AddCompensation(cleanup)

	docker := testApp.GetHttpApplication().Services.Docker

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
		projectNameSha := fmt.Sprintf("%.*x", 6, sha256.Sum256([]byte(projectName)))
		deploymentConfigPath := filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("buang-%v.yaml", projectNameSha))
		require.FileExists(t, deploymentConfigPath)

		cli, err := command.NewDockerCli(command.WithAPIClient(docker))
		require.NoError(t, err)
		err = cli.Initialize(&flags.ClientOptions{})
		require.NoError(t, err)

		svc, err := compose.NewComposeService(cli, compose.WithPrompt(compose.AlwaysOkPrompt()))
		require.NoError(t, err)

		composeProject, err := svc.LoadProject(ctx, api.ProjectLoadOptions{
			ConfigPaths: []string{filepath.Join(*deployment.Deployment.GetClonePath(), p.Project.GetComposePath())},
			ProjectName: projectNameSha,
		})
		require.NoError(t, err)

		require.Len(t, composeProject.AllServices(), EXPECTED_SERVICES)

		filters := filters.NewArgs()
		filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectNameSha))
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
			require.Equal(t, *envs["BUANG_PROJECT_NAME"], projectNameSha)

			aliases := []string{}
			for _, n := range i.NetworkSettings.Networks {
				aliases = append(aliases, n.Aliases...)
			}

			require.True(t, slices.ContainsFunc(aliases, func(alias string) bool {
				pattern := fmt.Sprintf("buang-%v-.+", projectNameSha)

				matched, err := regexp.MatchString(pattern, alias)
				if err != nil {
					return false
				}

				return matched
			}))
		}
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
		require.NoError(t, err)

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
	traefikId, traefikPort := createTraefik()
	getDeploymentLogs(projectId, deploymentId)
	assertDeployment(projectId, deploymentId)
	assertProxy(projectId, deploymentId, traefikId, traefikPort)
	buangDeployment(projectId, deploymentId)
}
