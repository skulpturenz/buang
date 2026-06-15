package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type ListDeploymentsParams struct {
	ProjectID int64
	Limit     int32
	Page      int32
}

type ListDeploymentsResult struct {
	Deployments []interfaces.Deployment
}

func (d ListDeploymentsParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ListDeploymentsResult, error) {
	assert.True(d.Limit > 0, "invalid limit")
	assert.True(d.Page > 0, "invalid page")

	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectDeploymentsDesc(ctx, interfaces.SelectDeploymentsDescParams(d))
	if err != nil {
		return nil, err
	}

	ret := ListDeploymentsResult{
		Deployments: result,
	}

	return &ret, nil
}
