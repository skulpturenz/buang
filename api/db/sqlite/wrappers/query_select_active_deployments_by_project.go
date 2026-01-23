package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectActiveDeploymentsByProject(ctx context.Context, projectId int64) ([]interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectActiveDeploymentsByProject(ctx, projectId)

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
