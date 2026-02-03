package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	"time"
)

type UpdateDeploymentParams struct {
	Url        *string
	Status     int16
	DeployedAt *time.Time
	ClonePath  *string
	ID         int64
	ProjectID  int64
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
