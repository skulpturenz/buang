package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) UpdateDeploymentStatus(ctx context.Context, arg interfaces.UpdateDeploymentStatusParams) (interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.UpdateDeploymentStatus(ctx, sqlitemodels.UpdateDeploymentStatusParams{
		Status: arg.Status,
		ID:     arg.ID,
	})

	return Deployment(ret), err
}
