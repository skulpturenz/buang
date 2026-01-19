package activities

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"skulpture/buang/components/docker"
	"skulpture/buang/components/projects"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
)

type BuangDeployment app.ApplicationServices

type BuangDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

type BuangDeploymentResult struct{}

func (bd *BuangDeployment) BuangDeployment(ctx context.Context, b BuangDeploymentParams) (*BuangDeploymentResult, error) {
	s := app.ApplicationServices(*bd)

	findProjectParams := projects.FindProjectByIdParams{
		ID: b.ProjectId,
	}

	p, err := findProjectParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	findDeploymentParams := deployments.FindDeploymentByIdParams{
		ID:        b.DeploymentId,
		ProjectId: p.Project.GetId(),
	}

	dply, err := findDeploymentParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	clonePath := *dply.Deployment.GetClonePath()
	deploymentConfigPath := filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("project-%v-deployment-%v.yaml", p.Project.GetId(), dply.Deployment.GetId()))
	sha := fmt.Sprintf("%.*s", 8, dply.Deployment.GetSha())
	projectName := fmt.Sprintf("%v_%v_%v_%v", p.Project.GetId(), dply.Deployment.GetId(), dply.Deployment.GetBranch(), sha)
	configPath := filepath.Join(clonePath, p.Project.GetComposePath())

	err = os.RemoveAll(deploymentConfigPath)
	if err != nil {
		return nil, err
	}

	downParams := docker.ComposeDownParams{
		ProjectName: projectName,
		ConfigPaths: []string{configPath},
	}

	_, err = downParams.Exec(ctx, &s)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("unable to compose down deployment %v for project %v: %v", dply.Deployment.GetId(), dply.Deployment.GetProjectId(), err.Error()))

		p := deployments.UpdateDeploymentParams{
			ID:         dply.Deployment.GetId(),
			Url:        dply.Deployment.GetUrl(),
			Status:     int16(enumsdeploymentstatus.Error),
			DeployedAt: dply.Deployment.GetDeployedAt(),
			ClonePath:  dply.Deployment.GetClonePath(),
		}

		_, err := p.Exec(ctx, &s)
		if err != nil {
			return nil, err
		}
	}

	err = os.RemoveAll(clonePath)
	if err != nil {
		return nil, err
	}

	buangParams := deployments.BuangDeploymentParams{
		ID: dply.Deployment.GetId(),
	}

	_, err = buangParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &BuangDeploymentResult{}, nil
}
