package dbosworkflows

import (
	"context"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

func BuangDeployment(ctx dbos.DBOSContext, p BuangDeploymentParams) (bool, error) {
	_, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment{}

			_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return &activities.BuangDeploymentResult{}, nil
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
