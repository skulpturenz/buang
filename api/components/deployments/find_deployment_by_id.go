package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindDeploymentByIdParams struct {
	ID        int64
	ProjectId int64
}

type FindDeploymentByIdResult struct {
	Deployment interfaces.Deployment
}

func (p FindDeploymentByIdParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindDeploymentByIdResult, error) {
	q := *s.Queries

	result, err := q.SelectDeployment(ctx, interfaces.SelectDeploymentParams{
		ID:        p.ID,
		ProjectID: p.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	ret := FindDeploymentByIdResult{
		Deployment: result,
	}

	return &ret, nil
}
