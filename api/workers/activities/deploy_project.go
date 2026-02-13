package activities

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"skulpture/buang/app"
	deploymentlogs "skulpture/buang/components/deployment_logs"
	"skulpture/buang/components/deployments"
	"skulpture/buang/components/docker"
	"skulpture/buang/components/projects"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
	"strings"
	"time"

	dynamic "github.com/traefik/traefik/v3/pkg/config/dynamic"
	"go.yaml.in/yaml/v3"
)

type DeployProject app.ApplicationServices

type DeployProjectParams struct {
	ProjectId    int64
	DeploymentId int64
	Dir          string
}

type DeployProjectResult struct{}

func (dp *DeployProject) DeployProject(ctx context.Context, d DeployProjectParams) (*DeployProjectResult, error) {
	DEPLOYED_AT := time.Now()

	s := app.ApplicationServices(*dp)

	err := os.MkdirAll(deployments.TRAEFIK_DYNAMIC_CONFIG, os.ModePerm)
	if err != nil {
		return nil, err
	}

	findProjectParams := projects.FindProjectByIdParams{
		ID: d.ProjectId,
	}

	p, err := findProjectParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	findDeploymentParams := deployments.FindDeploymentByIdParams{
		ID:        d.DeploymentId,
		ProjectId: p.Project.GetId(),
	}

	dply, err := findDeploymentParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	if dply.Deployment.GetStatus() != int16(enumsdeploymentstatus.Deploying) {
		return nil, fmt.Errorf("deployment %v for project %v is invalid", d.DeploymentId, d.ProjectId)
	}

	deploymentLogsParams := deploymentlogs.DeploymentLogWriterParams{
		ProjectID:    dply.Deployment.GetProjectId(),
		DeploymentID: dply.Deployment.GetId(),
	}

	_, err = deploymentLogsParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	deployingParams := deployments.UpdateDeploymentParams{
		ID:         dply.Deployment.GetId(),
		Status:     int16(enumsdeploymentstatus.Deploying),
		ProjectID:  p.Project.GetId(),
		Url:        nil,
		DeployedAt: &DEPLOYED_AT,
		ClonePath:  &d.Dir,
	}

	_, err = deployingParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	env, err := dply.Deployment.GetEnvVars()
	if err != nil {
		return nil, err
	}
	if env == nil {
		env = map[string]any{}
	}

	projectName := deployments.GetProjectName(deployments.GetProjectNameParams{
		ProjectId:    p.Project.GetId(),
		DeploymentId: dply.Deployment.GetId(),
		Branch:       dply.Deployment.GetBranch(),
		Sha:          dply.Deployment.GetSha(),
		DeployedAt:   DEPLOYED_AT,
	})
	url := fmt.Sprintf("/deployment/%v", projectName)

	traefikEntrypointRouter := projectName
	traefikEntrypointService := projectName
	env["BUANG_DEPLOYMENT_PATH"] = url
	env["BUANG_ENTRYPOINT_TRAEFIK_ROUTER"] = traefikEntrypointRouter
	env["BUANG_ENTRYPOINT_TRAEFIK_SERVICE"] = traefikEntrypointService,

	composePath := filepath.Join(d.Dir, p.Project.GetComposePath())
	expandedComposePath, _, err := docker.ExpandComposeYaml(composePath, env)
	if err != nil {
		return nil, err
	}

	// TODO: ideally we want to use the buffered writer so that we don't hit the DB all the time
	// but temporal's long polling is faster than our short polling in `poll_deployment_log`
	// so what happens is that the workflow unblocks and we return deployment logs for the deployment
	// but the compose logs are not included because it unblocks before `poll_deployment_logs` can
	// retrieve the updated logs
	// steps to reproduce:
	// - temporal with sqlite
	//   - has deployed check returns too early. after git clone is done, `HasDeployed` unblocks
	//     but docker compose logs are not written yet
	//     steps to reproduce:
	//       - create project
	//       - get deployment logs for project 1, deployment id 1 (this will be streaming)
	//       - create deployment (this will have id 1)
	//       - once the git clone is complete, the deployment logs returns but there are no logs
	//         from compose, only git
	//       - note: if we fetch it again the compose logs show so unblocking too early
	//               not so sure why though because the workflow isn't complete
	// fine with dbos
	// TODO: think this is something which won't occur when deployed because network latency will give us enough
	// time for everything to work correctly, in which case we can use the buffered writer but we just have to relax the
	// assertion around compose logs being included
	// writer := deploymentlogs.CreateBufferedWriter(deploymentLogsParams)
	// defer writer.Flush()
	upParams := docker.ComposeUpParams{
		ProjectName: projectName,
		ConfigPaths: []string{
			*expandedComposePath,
		},
		Environment: env,
		Writer:      deploymentLogsParams,
	}

	_, _, err = upParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	serviceEntrypoint := strings.Split(dply.Deployment.GetServiceEntrypoint(), ":")

	passHostHeader := true
	config := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				traefikEntrypointRouter: {
					EntryPoints: []string{"web", "websecure"},
					Rule:        fmt.Sprintf("PathPrefix(`%v`)", url),
					Service:     projectName,
					Middlewares: []string{fmt.Sprintf("%v-stripprefix", traefikEntrypointRouter)},
				},
			},
			Services: map[string]*dynamic.Service{
				traefikEntrypointService: {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: fmt.Sprintf("http://%v:%v", serviceEntrypoint[0], serviceEntrypoint[1])},
						},
						PassHostHeader: &passHostHeader,
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				fmt.Sprintf("%v-stripprefix", traefikEntrypointRouter): &dynamic.Middleware{
					StripPrefix: &dynamic.StripPrefix{
						Prefixes: []string{url},
					},
				},
			},
		},
	}

	yml, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}

	deploymentConfigPath := deployments.GetDeploymentPath(deployments.GetDeploymentPathParams{
		ProjectId:    p.Project.GetId(),
		DeploymentId: dply.Deployment.GetId(),
		Branch:       dply.Deployment.GetBranch(),
		Sha:          dply.Deployment.GetSha(),
		DeployedAt:   DEPLOYED_AT,
	})
	err = os.WriteFile(deploymentConfigPath, yml, 0644)
	if err != nil {
		return nil, err
	}

	deployedParams := deployments.UpdateDeploymentParams{
		ID:         dply.Deployment.GetId(),
		Url:        &url,
		DeployedAt: &DEPLOYED_AT,
		Status:     int16(enumsdeploymentstatus.Deployed),
		ClonePath:  dply.Deployment.GetClonePath(),
		ProjectID:  p.Project.GetId(),
	}

	_, err = deployedParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &DeployProjectResult{}, nil
}
