package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
)

type CreateDeploymentParams struct {
	ProjectID int64
	Sha       *string
}

type CreateDeploymentResult struct {
	Id int64
}

func (d CreateDeploymentParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CreateDeploymentResult, error) {
	q := *s.Queries

	result, err := q.CreateDeployment(ctx, interfaces.CreateDeploymentParams{
		ProjectID: d.ProjectID,
		Sha:       d.Sha,
		Status:    int16(enumsdeploymentstatus.New),
	})
	if err != nil {
		return nil, err
	}

	ret := CreateDeploymentResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
