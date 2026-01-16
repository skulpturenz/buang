package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectProjectsDesc(ctx context.Context, args interfaces.SelectProjectsDescParams) ([]interfaces.Project, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectProjectsDesc(ctx, pgmodels.SelectProjectsDescParams{
		Limit:   args.Limit,
		Column2: args.Page,
	})

	ret := []interfaces.Project{}
	for _, x := range res {
		ret = append(ret, Project(x))
	}

	return ret, err
}
