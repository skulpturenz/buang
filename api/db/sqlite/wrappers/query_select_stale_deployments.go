package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectStaleDeployments(ctx context.Context) ([]interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectStaleDeployments(ctx)

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
