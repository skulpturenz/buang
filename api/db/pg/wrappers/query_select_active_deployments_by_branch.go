package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectActiveDeploymentsByBranch(ctx context.Context, arg interfaces.SelectActiveDeploymentsByBranchParams) ([]interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectActiveDeploymentsByBranch(ctx, pgmodels.SelectActiveDeploymentsByBranchParams(arg))

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
