package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) CreateProject(ctx context.Context, arg interfaces.CreateProjectParams) (interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.CreateProject(ctx, sqlitemodels.CreateProjectParams{
		Repository:    arg.Repository,
		RequiresAuthn: arg.RequiresAuthn,
		Username:      arg.Username,
		Password:      arg.Password,
		ComposePath:   arg.ComposePath,
	})

	return Project(ret), err
}
