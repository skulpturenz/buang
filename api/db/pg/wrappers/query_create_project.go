package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) CreateProject(ctx context.Context, arg interfaces.CreateProjectParams) (interfaces.Project, error) {
	m := pgmodels.Queries(q)
	ret, err := m.CreateProject(ctx, pgmodels.CreateProjectParams(arg))

	return Project(ret), err
}
