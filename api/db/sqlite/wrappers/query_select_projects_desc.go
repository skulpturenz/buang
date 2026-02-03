package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectProjectsDesc(ctx context.Context, args interfaces.SelectProjectsDescParams) ([]interfaces.Project, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectProjectsDesc(ctx, sqlitemodels.SelectProjectsDescParams{
		Limit: int64(args.Limit),
		Page:  args.Page,
	})

	ret := []interfaces.Project{}
	for _, x := range res {
		ret = append(ret, Project(x))
	}

	return ret, err
}
