package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindProjectByRepositoryParams struct {
	Repository string
}

type FindProjectByRepositoryResult struct {
	Project interfaces.Project
}

func (p FindProjectByRepositoryParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindProjectByRepositoryResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectProjectByRepository(ctx, p.Repository)
	if err != nil {
		return nil, err
	}

	assert.True(result.GetId() > 0, "invalid result")

	ret := FindProjectByRepositoryResult{
		Project: result,
	}

	return &ret, nil
}
