package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectStaleDeployments(ctx context.Context) ([]interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectStaleDeployments(ctx)

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
