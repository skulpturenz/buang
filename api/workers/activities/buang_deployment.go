package activities

import (
	"context"
	"errors"
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

	clonePathPtr := dply.Deployment.GetClonePath()
	if clonePathPtr == nil {
		buangParams := deployments.BuangDeploymentParams{
			ID:        dply.Deployment.GetId(),
			ProjectID: dply.Deployment.GetProjectId(),
		}

		_, err = buangParams.Exec(ctx, &s)
		if err != nil {
			return nil, err
		}

		return &BuangDeploymentResult{}, nil
	}

	clonePath := *clonePathPtr
	projectName := deployments.GetProjectName(deployments.GetProjectNameParams{
		ProjectId:    p.Project.GetId(),
		DeploymentId: dply.Deployment.GetId(),
		Branch:       dply.Deployment.GetBranch(),
		Sha:          dply.Deployment.GetSha(),
		DeployedAt:   *dply.Deployment.GetDeployedAt(),
	})
	deploymentConfigPath := deployments.GetDeploymentPath(deployments.GetDeploymentPathParams{
		ProjectId:    p.Project.GetId(),
		DeploymentId: dply.Deployment.GetId(),
		Branch:       dply.Deployment.GetBranch(),
		Sha:          dply.Deployment.GetSha(),
		DeployedAt:   *dply.Deployment.GetDeployedAt(),
	})

	expandedPath := filepath.Join(clonePath, fmt.Sprintf("expanded-%v", filepath.Base(p.Project.GetComposePath())))
	_, err = os.Stat(expandedPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var configPath string
	if errors.Is(err, os.ErrNotExist) {
		configPath = filepath.Join(clonePath, p.Project.GetComposePath())
	} else {
		configPath = expandedPath
	}

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
			ProjectID:  dply.Deployment.GetProjectId(),
		}

		_, updateDeploymentErr := p.Exec(ctx, &s)
		if updateDeploymentErr != nil {
			return nil, errors.Join(err, updateDeploymentErr)
		}

		return nil, err
	}

	err = os.RemoveAll(clonePath)
	if err != nil {
		return nil, err
	}

	buangParams := deployments.BuangDeploymentParams{
		ID:        dply.Deployment.GetId(),
		ProjectID: dply.Deployment.GetProjectId(),
	}

	_, err = buangParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &BuangDeploymentResult{}, nil
}
