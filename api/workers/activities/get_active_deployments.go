package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
)

type GetActiveDeploymentIds app.ApplicationServices

type GetActiveDeploymentIdsParams struct {
	ProjectId int64
}

type GetActiveDeploymentIdsResult struct {
	DeploymentIds []int64
}

func (ad *GetActiveDeploymentIds) GetActiveDeploymentIds(ctx context.Context, p GetActiveDeploymentIdsParams) (*GetActiveDeploymentIdsResult, error) {
	s := app.ApplicationServices(*ad)

	params := deployments.FindActiveDeploymentsParams{
		ProjectId: p.ProjectId,
	}

	deployments, err := params.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	ids := []int64{}
	for _, d := range deployments.Deployments {
		ids = append(ids, d.GetId())
	}

	return &GetActiveDeploymentIdsResult{DeploymentIds: ids}, nil
}
