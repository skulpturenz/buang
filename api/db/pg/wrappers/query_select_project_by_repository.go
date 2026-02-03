package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectProjectByRepository(ctx context.Context, repository string) (interfaces.Project, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectProjectByRepository(ctx, repository)

	ret := Project(res)

	return ret, err
}
