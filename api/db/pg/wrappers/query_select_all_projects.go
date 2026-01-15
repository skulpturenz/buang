package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pg_models "skulpture/buang/db/pg/out"
)

type Queries pg_models.Queries

func (q Queries) SelectAllProjectsDesc(ctx context.Context) ([]interfaces.Project, error) {
	res, err := q.SelectAllProjectsDesc(ctx)

	ret := []interfaces.Project{}
	for _, x := range res {
		ret = append(ret, interfaces.Project(x))
	}

	return ret, err
}
