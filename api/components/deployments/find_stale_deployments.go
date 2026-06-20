package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindStaleDeploymentsParams struct {
	ID        int64
	ProjectId int64
}

type FindStaleDeploymentsResult struct {
	Deployments []interfaces.Deployment
}

func (p FindStaleDeploymentsParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindStaleDeploymentsResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectStaleDeployments(ctx)
	if err != nil {
		return nil, err
	}

	ret := FindStaleDeploymentsResult{
		Deployments: result,
	}

	return &ret, nil
}
