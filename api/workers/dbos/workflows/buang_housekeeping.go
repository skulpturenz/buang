package dbosworkflows

import (
	"context"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangHousekeepingParams struct{}

func BuangHousekeeping(ctx dbos.DBOSContext, _ BuangHousekeepingParams) (bool, error) {
	staleDeployments, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetStaleDeploymentsResult, error) {
			getStaleDeployments := activities.GetStaleDeployments{}

			res, err := getStaleDeployments.GetStaleDeployments(ctx, activities.GetStaleDeploymentsParams{})
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	)
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment{}

			for _, d := range staleDeployments.Deployments {
				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    d.ProjectID,
					DeploymentId: d.DeploymentID,
				})
				if err != nil {
					return nil, err
				}
			}

			return &activities.BuangDeploymentResult{}, nil
		},
	)
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.PruneResult, error) {
			prune := activities.Prune{}

			res, err := prune.Prune(ctx)
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	)

	return true, nil
}
