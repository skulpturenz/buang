package temporalworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangDeployment(ctx workflow.Context, projectId int64, deploymentId int64) error {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang deployment panic: %v", r))
		}
	}()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	err := workflow.ExecuteActivity(ctx,
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
