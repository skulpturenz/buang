package dbosworkflows

import (
	"context"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type DeployParams struct {
	ProjectId    int64
	DeploymentId int64
}

func Deploy(ctx dbos.DBOSContext, p DeployParams) (bool, error) {
	deploymentBranch, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetDeploymentBranchResult, error) {
			getDeploymentBranch := activities.GetDeploymentBranch{}
			res, err := getDeploymentBranch.GetDeploymentBranch(ctx, activities.GetDeploymentBranchParams{
				ProjectId: p.ProjectId,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepName("getDeploymentBranch"))
	if err != nil {
		return false, err
	}

	activeDeploymentIds, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetActiveDeploymentIdsResult, error) {
			getActiveDeploymentIds := activities.ActiveDeploymentIds{}
			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    deploymentBranch.Branch,
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

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.CreateDynamicConfigDirResult, error) {
			createDynamicConfigDir := activities.CreateDynamicConfigDir{}

			err := createDynamicConfigDir.CreateDynamicConfigDir(ctx, activities.CreateDynamicConfigDirParams{})
			if err != nil {
				return nil, err
			}

			return &activities.CreateDynamicConfigDirResult{}, nil
		},
	)

	cloneDeployment, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.CloneDeploymentResult, error) {
			cloneDeployment := activities.CloneDeployment{}

			res, err := cloneDeployment.CloneDeployment(ctx, activities.CloneDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	)

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.DeployProjectResult, error) {
			deployProject := activities.DeployProject{}

			res, err := deployProject.DeployProject(ctx, activities.DeployProjectParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
				Dir:          cloneDeployment.Dir,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	)

	return true, nil
}
