package dbosworkflows

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangDeployment app.ApplicationServices

type BuangDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

func (bd BuangDeployment) BuangDeployment(ctx dbos.DBOSContext, p BuangDeploymentParams) (bool, error) {
	s := app.ApplicationServices(bd)

	_, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment(s)

			_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return &activities.BuangDeploymentResult{}, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	return true, nil
}
