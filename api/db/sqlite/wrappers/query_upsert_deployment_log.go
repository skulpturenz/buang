package wrappers

import (
	"context"
	"skulpture/buang/db/interfaces"
	sqlite "skulpture/buang/db/sqlite/out"
)

func (q Queries) UpsertDeploymentLog(ctx context.Context, p interfaces.UpsertDeploymentLogParams) (interfaces.DeploymentLog, error) {
	m := sqlite.Queries(q)
	ret, err := m.UpsertDeploymentLog(ctx, sqlite.UpsertDeploymentLogParams{
		ProjectId:    p.ProjectID,
		DeploymentId: p.DeploymentID,
		Log:          p.Log,
	})
	if err != nil {
		return nil, err
	}

	return DeploymentLog(ret), nil
}
