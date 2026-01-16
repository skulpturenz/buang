package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindProjectByIdParams struct {
	ID int64
}

type FindProjectByIdResult struct {
	Project interfaces.Project
}

func (p FindProjectByIdParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindProjectByIdResult, error) {
	q := *s.Queries

	result, err := q.SelectProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	ret := FindProjectByIdResult{
		Project: result,
	}

	return &ret, nil
}
