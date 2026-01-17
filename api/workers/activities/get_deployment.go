package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"skulpture/buang/db/interfaces"
)

type GetDeployment app.ApplicationServices

type GetDeploymentParams struct {
	ProjectId int64
	ID        int64
}

type GetDeploymentResult struct {
	Deployment interfaces.Deployment
}

func (ad *GetDeployment) GetDeployment(ctx context.Context, p GetDeploymentParams) (*GetDeploymentResult, error) {
	s := app.ApplicationServices(*ad)

	params := deployments.FindDeploymentByIdParams{
		ID:        p.ID,
		ProjectId: p.ProjectId,
	}

	deployment, err := params.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &GetDeploymentResult{Deployment: deployment.Deployment}, nil
}
