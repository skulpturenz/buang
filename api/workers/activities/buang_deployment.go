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

	clonePath := ""
	deploymentConfigPath := filepath.Join(TRAEFIK_DYNAMIC_CONFIG, fmt.Sprintf("project-%v-deployment-%v", p.Project.GetId(), dply.Deployment.GetId()))
	projectName := fmt.Sprintf("%v-%v", p.Project.GetId(), dply.Deployment.GetSha())
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
		return nil, err
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
