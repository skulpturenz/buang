package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) UpdateDeployment(ctx context.Context, arg interfaces.UpdateDeploymentParams) (interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	ret, err := m.UpdateDeployment(ctx, sqlitemodels.UpdateDeploymentParams{
		ID:         arg.ID,
		Url:        arg.Url,
		Status:     arg.Status,
		DeployedAt: arg.DeployedAt,
		ClonePath:  arg.ClonePath,
		ProjectId:  arg.ProjectID,
	})

	return Deployment(ret), err
}
