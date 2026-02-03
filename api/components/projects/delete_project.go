package projects

import (
	"context"
	"skulpture/buang/app"

	"github.com/negrel/assert"
)

type DeleteProjectParams struct {
	Id int64
}

type DeleteProjectResult struct{}

func (p DeleteProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*DeleteProjectResult, error) {
	q := *s.Queries

	result, err := q.DeleteProject(ctx, p.Id)
	if err != nil {
		return nil, err
	}

	assert.True(result.GetId() > 0, "invalid result")

	ret := DeleteProjectResult{}

	return &ret, nil
}
