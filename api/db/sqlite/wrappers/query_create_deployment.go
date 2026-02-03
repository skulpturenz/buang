package wrappers

import (
	"context"
	"encoding/json"
	"skulpture/buang/db/interfaces"
	sqlitemodels "skulpture/buang/db/sqlite/out"
)

func (q Queries) CreateDeployment(ctx context.Context, arg interfaces.CreateDeploymentParams) (interfaces.Deployment, error) {
	m := sqlitemodels.Queries(q)

	envVars, err := json.Marshal(arg.EnvVars)
	if err != nil {
		return nil, err
	}

	ret, err := m.CreateDeployment(ctx, sqlitemodels.CreateDeploymentParams{
		ProjectId:         arg.ProjectID,
		Sha:               arg.Sha,
		Status:            arg.Status,
		ServiceEntrypoint: arg.ServiceEntrypoint,
		Branch:            arg.Branch,
		EnvVars:           envVars,
	})

	return Deployment(ret), err
}
