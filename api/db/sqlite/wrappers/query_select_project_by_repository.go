package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectProjectByRepository(ctx context.Context, repository string) (interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectProjectByRepository(ctx, repository)

	ret := Project(res)

	return ret, err
}
