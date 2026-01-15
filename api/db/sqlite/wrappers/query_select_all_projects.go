package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlite_models "skulpture/buang/db/sqlite/out"
)

type Queries sqlite_models.Queries

func (q Queries) SelectAllProjectsDesc(ctx context.Context) ([]interfaces.Project, error) {
	res, err := q.SelectAllProjectsDesc(ctx)

	ret := []interfaces.Project{}
	for _, x := range res {
		ret = append(ret, interfaces.Project(x))
	}

	return ret, err
}
