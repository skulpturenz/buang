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

	deploymentLogsParams := deploymentlogs.DeploymentLogWriterParams{
		ProjectID:    dply.Deployment.GetProjectId(),
		DeploymentID: dply.Deployment.GetId(),
	}

	_, err = deploymentLogsParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	deployingParams := deployments.UpdateDeploymentParams{
		ID:        dply.Deployment.GetId(),
		Status:    int16(enumsdeploymentstatus.Deploying),
		ProjectID: p.Project.GetId(),
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

	sha := fmt.Sprintf("%.*s", 8, dply.Deployment.GetSha())
	projectName := fmt.Sprintf("%v_%v_%v_%v", p.Project.GetId(), dply.Deployment.GetId(), dply.Deployment.GetBranch(), sha)
	url := fmt.Sprintf("/deployment/%v", projectName)

	env["BUANG_DEPLOYMENT_PATH"] = url

	writer := deploymentlogs.CreateBufferedWriter(deploymentLogsParams)
	upParams := docker.ComposeUpParams{
		ProjectName: projectName,
		ConfigPaths: []string{
			filepath.Join(d.Dir, p.Project.GetComposePath()),
		},
		Environment: env,
		Writer:      writer,
	}
	defer writer.Flush()

	_, _, err = upParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	serviceEntrypoint := strings.Split(dply.Deployment.GetServiceEntrypoint(), ":")

	passHostHeader := true
	config := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				projectName: {
					EntryPoints: []string{"web"},
					Rule:        fmt.Sprintf("PathPrefix(`%v`)", url),
					Service:     projectName,
					Middlewares: []string{fmt.Sprintf("%v-stripprefix", projectName)},
				},
			},
			Services: map[string]*dynamic.Service{
				projectName: {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{URL: fmt.Sprintf("http://%v:%v", serviceEntrypoint[0], serviceEntrypoint[1])},
						},
						PassHostHeader: &passHostHeader,
						// TODO: i'm not sure but i think once the initial request to the deployment is made, if we have sticky cookies
						// enabled, then every subsequent request should get forwarded to the correct service even if we omit the deployment path
						// since the domain stays the same, only the path changes
						// need to check
						Sticky: &dynamic.Sticky{
							Cookie: &dynamic.Cookie{
								Name:     projectName,
								HTTPOnly: true,
							},
						},
					},
				},
			},
			Middlewares: map[string]*dynamic.Middleware{
				fmt.Sprintf("%v-stripprefix", projectName): &dynamic.Middleware{
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

	err = os.WriteFile(filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("project-%v-deployment-%v.yaml", p.Project.GetId(), dply.Deployment.GetId())), yml, 0644)
	if err != nil {
		return nil, err
	}

	deployedAt := time.Now()
	deployedParams := deployments.UpdateDeploymentParams{
		ID:         dply.Deployment.GetId(),
		Url:        &url,
		DeployedAt: &deployedAt,
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
