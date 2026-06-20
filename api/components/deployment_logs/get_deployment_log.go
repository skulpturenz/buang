package deploymentlogs

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"

	"github.com/negrel/assert"
)

type GetDeploymentLogParams struct {
	ProjectId    int64
	DeploymentId int64
}

type GetDeploymentLogResult struct {
	Logs *string
}

func (p GetDeploymentLogParams) Exec(ctx context.Context, s *app.ApplicationServices) (*GetDeploymentLogResult, error) {
	q, ok := s.GetQueries()
	assert.True(ok, "queries service not found")

	res, err := q.SelectDeploymentLog(ctx, interfaces.SelectDeploymentLogParams{
		ProjectID:    p.ProjectId,
		DeploymentID: p.DeploymentId,
	})
	if err != nil {
		return nil, err
	}

	ret := GetDeploymentLogResult{
		Logs: res.GetLog(),
	}

	return &ret, nil
}
