package wrappers

import (
	"context"
	"encoding/json"
	"skulpture/buang/db/interfaces"
	pgmodels "skulpture/buang/db/pg/out"
)

func (q Queries) CreateDeployment(ctx context.Context, arg interfaces.CreateDeploymentParams) (interfaces.Deployment, error) {
	m := pgmodels.Queries(q)

	envVars, err := json.Marshal(arg.EnvVars)
	if err != nil {
		return nil, err
	}

	ret, err := m.CreateDeployment(ctx, pgmodels.CreateDeploymentParams{
		ProjectID:         arg.ProjectID,
		Sha:               arg.Sha,
		Status:            arg.Status,
		ServiceEntrypoint: arg.ServiceEntrypoint,
		Branch:            arg.Branch,
		EnvVars:           envVars,
	})

	return Deployment(ret), err
}
