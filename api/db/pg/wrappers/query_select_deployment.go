package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectDeployment(ctx context.Context, args interfaces.SelectDeploymentParams) (interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectDeployment(ctx, pgmodels.SelectDeploymentParams(args))

	ret := Deployment(res)

	return ret, err
}
