package projects

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type UpdateProjectParams struct {
	Repository    string
	RequiresAuthn bool
	Username      *string
	Password      *string
	ComposePath   string
	ProjectID     int64
}

type UpdateProjectResult struct {
	Id int64
}

func (p UpdateProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*UpdateProjectResult, error) {
	q := *s.Queries

	result, err := q.UpdateProject(ctx, interfaces.UpdateProjectParams(p))
	if err != nil {
		return nil, err
	}

	assert.True(result.GetId() > 0, "invalid result")

	ret := UpdateProjectResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
