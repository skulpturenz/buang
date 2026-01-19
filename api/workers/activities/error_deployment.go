package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
)

type ErrorDeployment app.ApplicationServices

type ErrorDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

type ErrorDeploymentResult struct{}

func (bd *ErrorDeployment) ErrorDeployment(ctx context.Context, b ErrorDeploymentParams) (*ErrorDeploymentResult, error) {
	s := app.ApplicationServices(*bd)

	findDeploymentParams := deployments.FindDeploymentByIdParams{
		ID:        b.DeploymentId,
		ProjectId: b.ProjectId,
	}

	dply, err := findDeploymentParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	updateParams := deployments.UpdateDeploymentParams{
		ID:         dply.Deployment.GetId(),
		Url:        dply.Deployment.GetUrl(),
		Status:     int16(enumsdeploymentstatus.Error),
		DeployedAt: dply.Deployment.GetDeployedAt(),
		ClonePath:  dply.Deployment.GetClonePath(),
	}

	_, err = updateParams.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &ErrorDeploymentResult{}, nil
}
