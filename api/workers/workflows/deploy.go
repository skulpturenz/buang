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

	var getActiveDeployments *activities.GetActiveDeploymentIds
	var getActiveDeploymentsResult activities.GetActiveDeploymentIdsResult

	err := workflow.ExecuteActivity(ctx, getActiveDeployments.GetActiveDeploymentIds, activities.GetActiveDeploymentIdsParams{ProjectId: projectId}).Get(ctx, &getActiveDeploymentsResult)
	if err != nil {
		return err
	}

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	for _, id := range getActiveDeploymentsResult.DeploymentIds {
		err := workflow.ExecuteActivity(ctx, buangDeployment.BuangDeployment, activities.BuangDeploymentParams{ProjectId: projectId, DeploymentId: id}).Get(ctx, &buangDeploymentResult)

		if err != nil {
			return err
		}
	}

	var createDynamicConfigDir *activities.CreateDynamicConfigDir
	var createDynamicConfigDirResult activities.CloneDeploymentResult

	err = workflow.ExecuteActivity(ctx, createDynamicConfigDir.CreateDynamicConfigDir, activities.CreateDynamicConfigDirParams{}).Get(ctx, &createDynamicConfigDirResult)
	if err != nil {
		return err
	}

	var cloneDeploymentParams *activities.CloneDeployment
	var cloneDeploymentResult activities.CloneDeploymentResult

	err = workflow.ExecuteActivity(ctx, cloneDeploymentParams.CloneDeployment, activities.CloneDeploymentParams{ProjectId: projectId, DeploymentId: deploymentId}).Get(ctx, &cloneDeploymentResult)
	if err != nil {
		return err
	}

	var deployProject *activities.DeployProject
	var deployProjectResult activities.DeployProjectResult

	err = workflow.ExecuteActivity(ctx, deployProject.DeployProject, activities.DeployProjectParams{ProjectId: projectId, DeploymentId: deploymentId}).Get(ctx, &deployProjectResult)
	if err != nil {
		return err
	}

	return nil
}
