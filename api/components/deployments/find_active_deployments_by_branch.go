package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindActiveDeploymentsByBranchParams struct {
	ProjectID int64
	Branch    string
}

type FindActiveDeploymentsByBranchResult struct {
	Deployments []interfaces.Deployment
}

func (p FindActiveDeploymentsByBranchParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindActiveDeploymentsByBranchResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectActiveDeploymentsByBranch(ctx, interfaces.SelectActiveDeploymentsByBranchParams(p))
	if err != nil {
		return nil, err
	}

	ret := FindActiveDeploymentsByBranchResult{
		Deployments: result,
	}

	return &ret, nil
}
