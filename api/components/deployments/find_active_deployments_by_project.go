package deployments

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type FindActiveDeploymentsByProjectParams struct {
	ProjectID int64
}

type FindActiveDeploymentsByProjectResult struct {
	Deployments []interfaces.Deployment
}

func (p FindActiveDeploymentsByProjectParams) Exec(ctx context.Context, s *app.ApplicationServices) (*FindActiveDeploymentsByProjectResult, error) {
	q := *s.Queries

	result, err := q.SelectActiveDeploymentsByProject(ctx, p.ProjectID)
	if err != nil {
		return nil, err
	}

	ret := FindActiveDeploymentsByProjectResult{
		Deployments: result,
	}

	return &ret, nil
}
