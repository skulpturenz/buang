package deployments

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	dberrors "skulpture/buang/db/db_errors"
	"skulpture/buang/db/interfaces"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"

	"github.com/negrel/assert"
)

type CreateDeploymentParams struct {
	ProjectID         int64
	Branch            string
	Sha               string
	ServiceEntrypoint string
	Env               map[string]any
}

type CreateDeploymentResult struct {
	Id int64
}

func (d CreateDeploymentParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CreateDeploymentResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.CreateDeployment(ctx, interfaces.CreateDeploymentParams{
		ProjectID:         d.ProjectID,
		Branch:            d.Branch,
		Sha:               d.Sha,
		Status:            int16(enumsdeploymentstatus.New),
		ServiceEntrypoint: d.ServiceEntrypoint,
		EnvVars:           d.Env,
	})
	if dberrors.IsNoRows(err) {
		return nil, fmt.Errorf("deployment for this branch in progress, try again later")
	}
	if err != nil {
		return nil, err
	}

	ret := CreateDeploymentResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
