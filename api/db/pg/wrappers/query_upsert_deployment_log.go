package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) UpsertDeploymentLog(ctx context.Context, p interfaces.UpsertDeploymentLogParams) (interfaces.DeploymentLog, error) {
	m := pgmodels.Queries(q)
	ret, err := m.UpsertDeploymentLog(ctx, pgmodels.UpsertDeploymentLogParams{
		ProjectID: p.ProjectID,
		ID:        p.DeploymentID,
		Log:       p.Log,
	})
	if err != nil {
		return nil, err
	}

	return DeploymentLog(ret), nil
}
