package temporalworkflows

import (
	"errors"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangDeployment(ctx workflow.Context, projectId int64, deploymentId int64) (err error) {
	defer workers.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "BuangDeployment", err)

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	err = workflow.ExecuteActivity(ctx,
		buangDeployment.BuangDeployment,
		activities.BuangDeploymentParams{
			ProjectId:    projectId,
			DeploymentId: deploymentId,
		}).
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

	return nil
}
