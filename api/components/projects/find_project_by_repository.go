package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindProjectByRepositoryParams struct {
	Repository string
}

type FindProjectByRepositoryResult struct {
	Project interfaces.Project
}

func (p FindProjectByRepositoryParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindProjectByRepositoryResult, error) {
	q := *s.Queries

	result, err := q.SelectProjectByRepository(ctx, p.Repository)
	if err != nil {
		return nil, err
	}

	ret := FindProjectByRepositoryResult{
		Project: result,
	}

	return &ret, nil
}
