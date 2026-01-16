package workflows

import (
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Buang(ctx workflow.Context, projectId int64, deploymentId int64) error {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	err := workflow.ExecuteActivity(ctx, buangDeployment.Exec, activities.BuangDeploymentParams{ProjectId: projectId, DeploymentId: deploymentId}).Get(ctx, buangDeploymentResult)
	if err != nil {
		return err
	}

	return nil
}
