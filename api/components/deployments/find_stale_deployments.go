package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindStaleDeploymentsParams struct {
	ID        int64
	ProjectId int64
}

type FindStaleDeploymentsResult struct {
	Deployments []interfaces.Deployment
}

func (p FindStaleDeploymentsParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindStaleDeploymentsResult, error) {
	q := *s.Queries

	result, err := q.SelectStaleDeployments(ctx)
	if err != nil {
		return nil, err
	}

	ret := FindStaleDeploymentsResult{
		Deployments: result,
	}

	return &ret, nil
}
