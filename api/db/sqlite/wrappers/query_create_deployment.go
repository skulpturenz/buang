package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) CreateDeployment(ctx context.Context, arg interfaces.CreateDeploymentParams) (interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.CreateDeployment(ctx, sqlitemodels.CreateDeploymentParams{
		ProjectId:         arg.ProjectID,
		Sha:               arg.Sha,
		Status:            arg.Status,
		ServiceEntrypoint: arg.ServiceEntrypoint,
	})

	return Deployment(ret), err
}
