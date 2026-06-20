package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type ListProjectParams struct {
	Limit int32
	Page  int32
}

type ListProjectsResult struct {
	Projects []interfaces.Project
}

func (p ListProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ListProjectsResult, error) {
	assert.True(p.Limit > 0, "invalid limit")
	assert.True(p.Page > 0, "invalid page")

	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.SelectProjectsDesc(ctx, interfaces.SelectProjectsDescParams(p))
	if err != nil {
		return nil, err
	}

	ret := ListProjectsResult{
		Projects: result,
	}

	return &ret, nil
}
