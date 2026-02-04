package dbosworkflows

import (
	"context"
	"errors"
	"skulpture/buang/app"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers/activities"
	workersshared "skulpture/buang/workers/shared"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type Deploy app.ApplicationServices

type DeployParams struct {
	ProjectId    int64
	DeploymentId int64
}

func (d Deploy) Deploy(ctx dbos.DBOSContext, p DeployParams) (res bool, err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "Deploy", err)

	s := app.ApplicationServices(d)

	errorDeployment := func() error {
		_, err := dbos.RunAsStep(ctx,
			func(ctx context.Context) (*activities.ErrorDeploymentResult, error) {
				errorDeployment := activities.ErrorDeployment(s)
				res, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

				return res, nil

			}, dbos.WithStepMaxRetries(10))

		return err
	}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

		return false, err
	}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

		return false, err
	}

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
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		errorDeploymentErr := errorDeployment()
		if errorDeploymentErr != nil {
			return false, errors.Join(err, errorDeploymentErr)
		}

		return false, err
	}

	return true, nil
}
