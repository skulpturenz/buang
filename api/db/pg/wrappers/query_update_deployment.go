package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) UpdateDeployment(ctx context.Context, arg interfaces.UpdateDeploymentParams) (interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	ret, err := m.UpdateDeployment(ctx, pgmodels.UpdateDeploymentParams(arg))

	return Deployment(ret), err
}
