package temporalworkflows

import (
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangBranch(ctx workflow.Context, projectId int64, branch string) error {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var activeDeployments *activities.ActiveDeploymentIds
	var getActiveDeploymentsResult activities.GetActiveDeploymentIdsResult

	err := workflow.
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
			return err
		}
	}

	return nil
}
