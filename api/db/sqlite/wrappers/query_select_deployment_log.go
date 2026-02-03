package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectDeploymentLog(ctx context.Context, args interfaces.SelectDeploymentLogParams) (interfaces.DeploymentLog, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectDeploymentLog(ctx, sqlitemodels.SelectDeploymentLogParams{
		ProjectId:    args.ProjectID,
		DeploymentId: args.DeploymentID,
	})

	ret := DeploymentLog(res)

	return ret, err
}
