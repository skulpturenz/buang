package projects

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
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
	q := *s.Queries

	activeProject, err := q.SelectProjectByRepository(ctx, p.Repository)
	if err != nil {
		return nil, err
	}

	if !activeProject.GetDeleted() {
		return nil, fmt.Errorf("repository %v already has an active project with id %v", activeProject.GetRepository(), activeProject.GetId())
	}

	result, err := q.CreateProject(ctx, interfaces.CreateProjectParams(p))
	if err != nil {
		return nil, err
	}

	ret := CreateProjectResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
