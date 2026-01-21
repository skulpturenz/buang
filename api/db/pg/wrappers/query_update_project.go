package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) UpdateProject(ctx context.Context, arg interfaces.UpdateProjectParams) (interfaces.Project, error) {
	m := pgmodels.Queries(q)
	ret, err := m.UpdateProject(ctx, pgmodels.UpdateProjectParams(arg))

	return Project(ret), err
}
