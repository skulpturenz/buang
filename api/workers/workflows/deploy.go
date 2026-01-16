package workflows

import (
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Deploy(ctx workflow.Context, projectId int64, deploymentId int64) error {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var createDynamicConfigDir *activities.CreateDynamicConfigDir
	var createDynamicConfigDirResult activities.CloneDeploymentResult

	err := workflow.ExecuteActivity(ctx, createDynamicConfigDir.Exec, activities.CreateDynamicConfigDirParams{}).Get(ctx, createDynamicConfigDirResult)
	if err != nil {
		return err
	}

	var cloneDeploymentParams *activities.CloneDeployment
	var cloneDeploymentResult activities.CloneDeploymentResult

	err = workflow.ExecuteActivity(ctx, cloneDeploymentParams.Exec, activities.CloneDeploymentParams{ProjectId: projectId, DeploymentId: deploymentId}).Get(ctx, cloneDeploymentResult)
	if err != nil {
		return err
	}

	var deployProject *activities.DeployProject
	var deployProjectResult activities.DeployProjectResult

	err = workflow.ExecuteActivity(ctx, deployProject.Exec, activities.DeployProjectParams{ProjectId: projectId, DeploymentId: deploymentId}).Get(ctx, deployProjectResult)
	if err != nil {
		return err
	}

	return nil
}
