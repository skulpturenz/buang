package temporalworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Deploy(ctx workflow.Context, projectId int64, deploymentId int64) (err error) {
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

	// the convoluted error handling is because if this fails somewhere
	// and the deployment is still marked as new or deploying then no other deployments for the branch can happen

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var getDeploymentBranch *activities.GetDeploymentBranch
	var getDeploymentBranchResult activities.GetDeploymentBranchResult

	err = workflow.
		ExecuteActivity(ctx,
			getDeploymentBranch.GetDeploymentBranch,
			activities.GetDeploymentBranchParams{
				ProjectId: projectId, ID: deploymentId,
			}).
		Get(ctx, &getDeploymentBranchResult)
	if err != nil {
		return err
	}

	var activeDeploymentIds *activities.ActiveDeploymentIds
	var getActiveDeploymentsResult activities.GetActiveDeploymentIdsResult

	err = workflow.
		ExecuteActivity(ctx,
			activeDeploymentIds.GetActiveDeploymentIds,
			activities.GetActiveDeploymentIdsParams{
				ProjectId: projectId,
				Branch:    getDeploymentBranchResult.Branch,
			}).
		Get(ctx, &getActiveDeploymentsResult)
	if err != nil {
		var errorDeployment *activities.ErrorDeployment
		var errorDeploymentResult activities.ErrorDeploymentResult

		errErrorDeployment := workflow.ExecuteActivity(ctx,
			errorDeployment.ErrorDeployment,
			activities.ErrorDeploymentParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
			}).
			Get(ctx, &errorDeploymentResult)
		if errErrorDeployment != nil {
			return errors.Join(err, errErrorDeployment)
		}

		return err
	}

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	for _, id := range getActiveDeploymentsResult.DeploymentIds {
		err := workflow.
			ExecuteActivity(ctx,
				buangDeployment.BuangDeployment,
				activities.BuangDeploymentParams{
					ProjectId:    projectId,
					DeploymentId: id,
				},
			).
			Get(ctx, &buangDeploymentResult)

		if err != nil {
			var errorDeployment *activities.ErrorDeployment
			var errorDeploymentResult activities.ErrorDeploymentResult

			errErrorDeployment := workflow.ExecuteActivity(ctx,
				errorDeployment.ErrorDeployment,
				activities.ErrorDeploymentParams{
					ProjectId:    projectId,
					DeploymentId: deploymentId,
				}).
				Get(ctx, &errorDeploymentResult)
			if errErrorDeployment != nil {
				return errors.Join(err, errErrorDeployment)
			}

			return err
		}
	}

	var createDynamicConfigDir *activities.CreateDynamicConfigDir
	var createDynamicConfigDirResult activities.CloneDeploymentResult

	err = workflow.
		ExecuteActivity(ctx,
			createDynamicConfigDir.CreateDynamicConfigDir,
			activities.CreateDynamicConfigDirParams{},
		).
		Get(ctx, &createDynamicConfigDirResult)
	if err != nil {
		var errorDeployment *activities.ErrorDeployment
		var errorDeploymentResult activities.ErrorDeploymentResult

		errErrorDeployment := workflow.ExecuteActivity(ctx,
			errorDeployment.ErrorDeployment,
			activities.ErrorDeploymentParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
			}).
			Get(ctx, &errorDeploymentResult)
		if errErrorDeployment != nil {
			return errors.Join(err, errErrorDeployment)
		}

		return err
	}

	var cloneDeploymentParams *activities.CloneDeployment
	var cloneDeploymentResult activities.CloneDeploymentResult

	err = workflow.
		ExecuteActivity(ctx,
			cloneDeploymentParams.CloneDeployment,
			activities.CloneDeploymentParams{
				ProjectId: projectId, DeploymentId: deploymentId,
			}).
		Get(ctx, &cloneDeploymentResult)
	if err != nil {
		var errorDeployment *activities.ErrorDeployment
		var errorDeploymentResult activities.ErrorDeploymentResult

		errErrorDeployment := workflow.ExecuteActivity(ctx,
			errorDeployment.ErrorDeployment,
			activities.ErrorDeploymentParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
			}).
			Get(ctx, &errorDeploymentResult)
		if errErrorDeployment != nil {
			return errors.Join(err, errErrorDeployment)
		}

		return err
	}

	var deployProject *activities.DeployProject
	var deployProjectResult activities.DeployProjectResult

	err = workflow.
		ExecuteActivity(ctx,
			deployProject.DeployProject,
			activities.DeployProjectParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
				Dir:          cloneDeploymentResult.Dir,
			}).
		Get(ctx, &deployProjectResult)
	if err != nil {
		var errorDeployment *activities.ErrorDeployment
		var errorDeploymentResult activities.ErrorDeploymentResult

		errErrorDeployment := workflow.ExecuteActivity(ctx,
			errorDeployment.ErrorDeployment,
			activities.ErrorDeploymentParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
			}).
			Get(ctx, &errorDeploymentResult)
		if errErrorDeployment != nil {
			return errors.Join(err, errErrorDeployment)
		}

		return err
	}

	return nil
}
