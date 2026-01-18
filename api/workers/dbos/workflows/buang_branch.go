package dbosworkflows

import (
	"context"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangBranchParams struct {
	ProjectId int64
	Branch    string
}

func BuangBranch(ctx dbos.DBOSContext, p BuangBranchParams) (bool, error) {
	activeDeploymentIds, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetActiveDeploymentIdsResult, error) {
			getActiveDeploymentIds := activities.ActiveDeploymentIds{}
			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    p.Branch,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		})
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment{}

			for _, id := range activeDeploymentIds.DeploymentIds {
				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: id,
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

	return true, nil
}
