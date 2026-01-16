package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectDeploymentsDesc(ctx context.Context, args interfaces.SelectDeploymentsDescParams) ([]interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectDeploymentsDesc(ctx, sqlitemodels.SelectDeploymentsDescParams{
		ProjectId: args.ProjectID,
		Limit:     int64(args.Limit),
		Page:      args.Page,
	})

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
