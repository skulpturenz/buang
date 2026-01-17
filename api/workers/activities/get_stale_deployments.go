package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
)

type GetStaleDeployments app.ApplicationServices

type GetStaleDeploymentsParams struct {
}

type GetStaleDeploymentsResult struct {
	Deployments []struct {
		ProjectID    int64
		DeploymentID int64
	}
}

func (sd *GetStaleDeployments) GetStaleDeployments(ctx context.Context, p GetStaleDeploymentsParams) (*GetStaleDeploymentsResult, error) {
	s := app.ApplicationServices(*sd)

	params := deployments.FindStaleDeploymentsParams{}

	deployments, err := params.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	ds := []struct {
		ProjectID    int64
		DeploymentID int64
	}{}
	for _, d := range deployments.Deployments {
		ds = append(ds, struct {
			ProjectID    int64
			DeploymentID int64
		}{
			ProjectID:    d.GetProjectId(),
			DeploymentID: d.GetId(),
		})
	}

	return &GetStaleDeploymentsResult{Deployments: ds}, nil
}
