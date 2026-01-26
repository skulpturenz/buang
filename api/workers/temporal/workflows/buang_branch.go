package temporalworkflows

import (
	"errors"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers/activities"
	workersshared "skulpture/buang/workers/shared"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangBranch(ctx workflow.Context, projectId int64, branch string) (err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "BuangBranch", err)

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var activeDeployments *activities.ActiveDeploymentIds
	var getActiveDeploymentsResult activities.GetActiveDeploymentIdsResult

	err = workflow.
		ExecuteActivity(ctx,
			activeDeployments.GetActiveDeploymentIds,
			activities.GetActiveDeploymentIdsParams{
				ProjectId: projectId,
				Branch:    branch,
			}).
		Get(ctx, &getActiveDeploymentsResult)
	if err != nil {
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
					DeploymentId: id,
				}).
				Get(ctx, &errorDeploymentResult)
			if errErrorDeployment != nil {
				return errors.Join(err, errErrorDeployment)
			}

			return err
		}
	}

	return nil
}
