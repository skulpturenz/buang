package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) SelectDeploymentsDesc(ctx context.Context, args interfaces.SelectDeploymentsDescParams) ([]interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	res, err := m.SelectDeploymentsDesc(ctx, pgmodels.SelectDeploymentsDescParams{
		ProjectID: args.ProjectID,
		Limit:     args.Limit,
		Column3:   args.Page,
	})

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
