package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) UpdateProject(ctx context.Context, arg interfaces.UpdateProjectParams) (interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.UpdateProject(ctx, sqlitemodels.UpdateProjectParams{
		Repository:    arg.Repository,
		RequiresAuthn: arg.RequiresAuthn,
		Username:      arg.Username,
		Password:      arg.Password,
		ComposePath:   arg.ComposePath,
		ProjectId:     arg.ProjectID,
	})

	return Project(ret), err
}
