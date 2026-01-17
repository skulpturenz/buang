package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindActiveDeploymentsParams struct {
	ProjectId int64
}

type FindActiveDeploymentsResult struct {
	Deployments []interfaces.Deployment
}

func (p FindActiveDeploymentsParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindActiveDeploymentsResult, error) {
	q := *s.Queries

	result, err := q.SelectActiveDeployments(ctx, p.ProjectId)
	if err != nil {
		return nil, err
	}

	ret := FindActiveDeploymentsResult{
		Deployments: result,
	}

	return &ret, nil
}
