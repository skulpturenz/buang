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

type BuangBranch app.ApplicationServices

type BuangBranchParams struct {
	ProjectId int64
	Branch    string
}

func (bb BuangBranch) BuangBranch(ctx dbos.DBOSContext, p BuangBranchParams) (res bool, err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "BuangBranch", err)

	s := app.ApplicationServices(bb)

	errorDeployment := func(deploymentId int) error {
		_, err := dbos.RunAsStep(ctx,
			func(ctx context.Context) (*activities.ErrorDeploymentResult, error) {
				errorDeployment := activities.ErrorDeployment(s)
				res, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: int64(deploymentId),
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

				return res, nil

			}, dbos.WithStepMaxRetries(3))

		return err
	}

	activeDeploymentIds, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetActiveDeploymentIdsResult, error) {
			getActiveDeploymentIds := activities.ActiveDeploymentIds(s)
			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    p.Branch,
			})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	for _, id := range activeDeploymentIds.DeploymentIds {
		_, err = dbos.RunAsStep(ctx,
			func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
				buangDeployment := activities.BuangDeployment(s)

				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: id,
				})
				if err != nil {
					return nil, err
				}

				return &activities.BuangDeploymentResult{}, nil
			}, dbos.WithStepMaxRetries(3))
		if err != nil {
			errDeploymentErr := errorDeployment(int(id))
			if errDeploymentErr != nil {
				return false, errors.Join(err, errDeploymentErr)
			}

			return false, err
		}
	}

	return true, nil
}
