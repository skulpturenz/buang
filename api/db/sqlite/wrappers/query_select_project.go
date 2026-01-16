package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectProject(ctx context.Context, id int64) (interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectProject(ctx, id)

	ret := Project(res)

	return ret, err
}
