package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindDeploymentByIdParams struct {
	ID        int64
	ProjectId int64
}

type FindDeploymentByIdResult struct {
	Deployment interfaces.Deployment
}

func (p FindDeploymentByIdParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindDeploymentByIdResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

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
