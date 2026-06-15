package projects

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	dberrors "skulpture/buang/db/db_errors"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type CreateProjectParams struct {
	Repository    string
	RequiresAuthn bool
	Username      *string
	Password      *string
	ComposePath   string
}

type CreateProjectResult struct {
	Id int64
}

func (p CreateProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CreateProjectResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	activeProject, err := q.SelectProjectByRepository(ctx, p.Repository)
	if err != nil && !dberrors.IsNoRows(err) {
		return nil, err
	}

	if err == nil && !activeProject.GetDeleted() {
		return nil, fmt.Errorf("repository %v already has an active project with id %v", activeProject.GetRepository(), activeProject.GetId())
	}

	result, err := q.CreateProject(ctx, interfaces.CreateProjectParams(p))
	if err != nil {
		return nil, err
	}

	assert.True(result.GetId() > 0, "invalid result")

	ret := CreateProjectResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
