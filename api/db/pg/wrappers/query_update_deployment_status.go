package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) UpdateDeploymentStatus(ctx context.Context, arg interfaces.UpdateDeploymentStatusParams) (interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	ret, err := m.UpdateDeploymentStatus(ctx, pgmodels.UpdateDeploymentStatusParams(arg))

	return Deployment(ret), err
}
