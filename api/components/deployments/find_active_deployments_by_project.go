package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindActiveDeploymentsByProjectParams struct {
	ProjectID int64
}

type FindActiveDeploymentsByProjectResult struct {
	Deployments []interfaces.Deployment
}

func (p FindActiveDeploymentsByProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindActiveDeploymentsByProjectResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectActiveDeploymentsByProject(ctx, p.ProjectID)
	if err != nil {
		return nil, err
	}

	ret := FindActiveDeploymentsByProjectResult{
		Deployments: result,
	}

	return &ret, nil
}
