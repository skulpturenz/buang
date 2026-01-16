package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type CreateDeploymentParams struct {
	ProjectID int64
	Sha       *string
	Status    int16
}

type CreateDeploymentResult struct {
	Id int64
}

func (d CreateDeploymentParams) Exec(ctx context.Context, s *app.ApplicationServices) (*CreateDeploymentResult, error) {
	q := *s.Queries

	result, err := q.CreateDeployment(ctx, interfaces.CreateDeploymentParams(d))
	if err != nil {
		return nil, err
	}

	ret := CreateDeploymentResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
