package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type FindProjectByIdParams struct {
	ID int64
}

type FindProjectByIdResult struct {
	Project interfaces.Project
}

func (p FindProjectByIdParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindProjectByIdResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	assert.True(result.GetId() > 0, "invalid result")

	ret := FindProjectByIdResult{
		Project: result,
	}

	return &ret, nil
}
