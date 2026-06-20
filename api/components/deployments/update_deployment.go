package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	"time"

	"github.com/negrel/assert"
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
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	result, err := q.UpdateDeployment(ctx, interfaces.UpdateDeploymentParams(d))
	if err != nil {
		return nil, err
	}

	ret := UpdateDeploymentResult{
		Id: result.GetId(),
	}

	return &ret, nil
}
