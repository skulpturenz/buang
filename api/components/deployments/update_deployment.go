package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	"time"
)

type UpdateDeploymentParams struct {
	ID         int64
	Url        *string
	Status     int16
	DeployedAt *time.Time
	ClonePath  *string
}

type UpdateDeploymentResult struct {
	Id int64
}

func (d UpdateDeploymentParams) Exec(ctx context.Context, s *app.ApplicationServices) (*UpdateDeploymentResult, error) {
	q := *s.Queries

	result, err := q.UpdateDeployment(ctx, interfaces.UpdateDeploymentParams(d))
	if err != nil {
		return nil, err
	}

	ret := UpdateDeploymentResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
