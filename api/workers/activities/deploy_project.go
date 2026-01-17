package activities

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"skulpture/buang/app"
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
	s := app.ApplicationServices(*dp)

	err := os.MkdirAll(TRAEFIK_DYNAMIC_CONFIG, os.ModePerm)
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

	deployingParams := deployments.UpdateDeploymentParams{
		ID:     dply.Deployment.GetId(),
		Status: int16(enumsdeploymentstatus.Deploying),
	}

	_, err = deployingParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	env, err := dply.Deployment.GetEnvVars()
	if err != nil {
		return nil, err
	}

	upParams := docker.ComposeUpParams{
		ProjectName: fmt.Sprintf("%v-%v", p.Project.GetId(), dply.Deployment.GetSha()),
		ConfigPaths: []string{
			filepath.Join(d.Dir, p.Project.GetComposePath()),
		},
		Environment: env,
	}

	_, _, err = upParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	entrypointName := fmt.Sprintf("%v-%v-%v", p.Project.GetId(), dply.Deployment.GetId(), dply.Deployment.GetSha())
	serviceEntrypoint := strings.Split(dply.Deployment.GetServiceEntrypoint(), ":")
	url := fmt.Sprintf("/deployment/%v", dply.Deployment.GetSha())
	config := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				entrypointName: {
					EntryPoints: []string{"http"},
					Rule:        fmt.Sprintf("/deployment/%v", dply.Deployment.GetSha()),
					Service:     dply.Deployment.GetSha(),
				},
			},
			Services: map[string]*dynamic.Service{
				dply.Deployment.GetSha(): {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: serviceEntrypoint[0], Port: serviceEntrypoint[1]},
						},
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				fmt.Sprintf("%v-stripprefix", entrypointName): &dynamic.Middleware{
					StripPrefix: &dynamic.StripPrefix{
						Prefixes: []string{fmt.Sprintf("/deployment/%v", dply.Deployment.GetSha())},
					},
				},
			},
		},
	}

	yml, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("project-%v-deployment-%v", p.Project.GetId(), dply.Deployment.GetId())), yml, 0644)
	if err != nil {
		return nil, err
	}

	deployedAt := time.Now()
	deployedParams := deployments.UpdateDeploymentParams{
		ID:         dply.Deployment.GetId(),
		Url:        &url,
		DeployedAt: &deployedAt,
		Status:     int16(enumsdeploymentstatus.Deployed),
	}

	_, err = deployedParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &DeployProjectResult{}, nil
}
