package wrappers

import (
	"context"
	interfaces "skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) SelectActiveDeploymentsByBranch(ctx context.Context, arg interfaces.SelectActiveDeploymentsByBranchParams) ([]interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)
	res, err := m.SelectActiveDeploymentsByBranch(ctx, sqlitemodels.SelectActiveDeploymentsByBranchParams{
		ProjectId: arg.ProjectID,
		Branch:    arg.Branch,
	})

	ret := []interfaces.Deployment{}
	for _, x := range res {
		ret = append(ret, Deployment(x))
	}

	return ret, err
}
