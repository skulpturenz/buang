package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
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
	q := *s.Queries

	result, err := q.SelectDeploymentsDesc(ctx, interfaces.SelectDeploymentsDescParams(d))
	if err != nil {
		return nil, err
	}

	ret := ListDeploymentsResult{
		Deployments: result,
	}

	return &ret, nil
}
