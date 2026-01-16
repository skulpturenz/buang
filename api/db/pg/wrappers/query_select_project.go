package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectProject(ctx context.Context, id int64) (interfaces.Project, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectProject(ctx, id)

	ret := Project(res)

	return ret, err
}
