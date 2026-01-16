package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) CreateDeployment(ctx context.Context, arg interfaces.CreateDeploymentParams) (interfaces.Deployment, error) {
	m := pgmodels.Queries(q)
	ret, err := m.CreateDeployment(ctx, pgmodels.CreateDeploymentParams(arg))

	return Deployment(ret), err
}
