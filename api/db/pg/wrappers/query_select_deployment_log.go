package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectDeploymentLog(ctx context.Context, args interfaces.SelectDeploymentLogParams) (interfaces.DeploymentLog, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectDeploymentLog(ctx, pgmodels.SelectDeploymentLogParams(args))

	ret := DeploymentLog(res)

	return ret, err
}
