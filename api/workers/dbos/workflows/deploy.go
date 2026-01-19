package dbosworkflows

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type Deploy app.ApplicationServices

type DeployParams struct {
	ProjectId    int64
	DeploymentId int64
}

func (d Deploy) Deploy(ctx dbos.DBOSContext, p DeployParams) (bool, error) {
	s := app.ApplicationServices(d)

	deploymentBranch, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetDeploymentBranchResult, error) {
			getDeploymentBranch := activities.GetDeploymentBranch(s)
			res, err := getDeploymentBranch.GetDeploymentBranch(ctx, activities.GetDeploymentBranchParams{
				ProjectId: p.ProjectId,
				ID:        p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	activeDeploymentIds, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetActiveDeploymentIdsResult, error) {
			getActiveDeploymentIds := activities.ActiveDeploymentIds(s)
			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    deploymentBranch.Branch,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment(s)

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
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.CreateDynamicConfigDirResult, error) {
			createDynamicConfigDir := activities.CreateDynamicConfigDir(s)

			err := createDynamicConfigDir.CreateDynamicConfigDir(ctx, activities.CreateDynamicConfigDirParams{})
			if err != nil {
				return nil, err
			}

			return &activities.CreateDynamicConfigDirResult{}, nil
		}, dbos.WithStepMaxRetries(3))

	cloneDeployment, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.CloneDeploymentResult, error) {
			cloneDeployment := activities.CloneDeployment(s)

			res, err := cloneDeployment.CloneDeployment(ctx, activities.CloneDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.DeployProjectResult, error) {
			deployProject := activities.DeployProject(s)

			res, err := deployProject.DeployProject(ctx, activities.DeployProjectParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
				Dir:          cloneDeployment.Dir,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))

	return true, nil
}
