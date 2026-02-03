package deploymentlogs

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type GetDeploymentLogParams struct {
	ProjectId    int64
	DeploymentId int64
}

type GetDeploymentLogResult struct {
	Logs *string
}

func (p GetDeploymentLogParams) Exec(ctx context.Context, s *app.ApplicationServices) (*GetDeploymentLogResult, error) {
	q := *s.Queries

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
