package activities

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
)

type GetDeploymentBranch app.ApplicationServices

type GetDeploymentBranchParams struct {
	ProjectId int64
	ID        int64
}

type GetDeploymentBranchResult struct {
	Branch string
}

func (ad *GetDeploymentBranch) GetDeploymentBranch(ctx context.Context, p GetDeploymentBranchParams) (*GetDeploymentBranchResult, error) {
	s := app.ApplicationServices(*ad)

	params := deployments.FindDeploymentByIdParams{
		ID:        p.ID,
		ProjectId: p.ProjectId,
	}

	deployment, err := params.Exec(ctx, &s)
	if err != nil {
		return nil, err
	}

	return &GetDeploymentBranchResult{Branch: deployment.Deployment.GetBranch()}, nil
}
