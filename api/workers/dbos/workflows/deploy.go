package dbosworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type Deploy app.ApplicationServices

type DeployParams struct {
	ProjectId    int64
	DeploymentId int64
}

func (d Deploy) Deploy(ctx dbos.DBOSContext, p DeployParams) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("deploy panic: %v", r))

			log := map[string]any{
				"stack": string(debug.Stack()),
			}

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  log,
			}
			p.Exec(context.Background())

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("deploy panic")
			}
		}
	}()

	s := app.ApplicationServices(d)

	// the convoluted error handling is because if this fails somewhere
	// and the deployment is still marked as new or deploying then no other deployments can happen

	deploymentBranch, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetDeploymentBranchResult, error) {
			getDeploymentBranch := activities.GetDeploymentBranch(s)
			errorDeployment := activities.ErrorDeployment(s)

			res, err := getDeploymentBranch.GetDeploymentBranch(ctx, activities.GetDeploymentBranchParams{
				ProjectId: p.ProjectId,
				ID:        p.DeploymentId,
			})
			if err != nil {
				_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

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
			errorDeployment := activities.ErrorDeployment(s)

			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    deploymentBranch.Branch,
			})
			if err != nil {
				_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

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
			errorDeployment := activities.ErrorDeployment(s)

			for _, id := range activeDeploymentIds.DeploymentIds {
				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: id,
				})
				if err != nil {
					_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
						ProjectId:    p.ProjectId,
						DeploymentId: p.DeploymentId,
					})
					if errErrorDeployment != nil {
						return nil, errors.Join(err, errErrorDeployment)
					}

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
			errorDeployment := activities.ErrorDeployment(s)

			err := createDynamicConfigDir.CreateDynamicConfigDir(ctx, activities.CreateDynamicConfigDirParams{})
			if err != nil {
				_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

				return nil, err
			}

			return &activities.CreateDynamicConfigDirResult{}, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	cloneDeployment, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.CloneDeploymentResult, error) {
			cloneDeployment := activities.CloneDeployment(s)
			errorDeployment := activities.ErrorDeployment(s)

			res, err := cloneDeployment.CloneDeployment(ctx, activities.CloneDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.DeployProjectResult, error) {
			deployProject := activities.DeployProject(s)
			errorDeployment := activities.ErrorDeployment(s)

			res, err := deployProject.DeployProject(ctx, activities.DeployProjectParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
				Dir:          cloneDeployment.Dir,
			})
			if err != nil {
				_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if errErrorDeployment != nil {
					return nil, errors.Join(err, errErrorDeployment)
				}

				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	return true, nil
}
