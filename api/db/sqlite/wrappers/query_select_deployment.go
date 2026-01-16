package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectDeployment(ctx context.Context, args interfaces.SelectDeploymentParams) (interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectDeployment(ctx, sqlitemodels.SelectDeploymentParams{
		ID:        args.ID,
		ProjectId: args.ProjectID,
	})

	ret := Deployment(res)

	return ret, err
}
