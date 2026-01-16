package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type ListProjectParams struct {
	Limit int32
	Page  int32
}

type ListProjectsResult struct {
	Projects []interfaces.Project
}

func (p ListProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*ListProjectsResult, error) {
	q := *s.Queries

	result, err := q.SelectProjectsDesc(ctx, interfaces.SelectProjectsDescParams(p))
	if err != nil {
		return nil, err
	}

	ret := ListProjectsResult{
		Projects: result,
	}

	return &ret, nil
}
