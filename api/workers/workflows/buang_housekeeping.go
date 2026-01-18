package workflows

import (
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CronResult struct {
	RunTime time.Time
}

func BuangHouskeeping(ctx workflow.Context) (*CronResult, error) {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	now := workflow.Now(ctx)

	var getStaleDeployments *activities.GetStaleDeployments
	var getStaleDeploymentsResult activities.GetStaleDeploymentsResult

	err := workflow.
		ExecuteActivity(ctx,
			getStaleDeployments.GetStaleDeployments).
		Get(ctx, &getStaleDeploymentsResult)
	if err != nil {
		return nil, err
	}

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	for _, d := range getStaleDeploymentsResult.Deployments {
		err := workflow.
			ExecuteActivity(ctx,
				buangDeployment.BuangDeployment,
				activities.BuangDeploymentParams{
					ProjectId:    d.ProjectID,
					DeploymentId: d.DeploymentID,
				},
			).
			Get(ctx, &buangDeploymentResult)

		if err != nil {
			return nil, err
		}
	}

	var prune *activities.Prune
	var pruneResult activities.PruneResult

	workflow.
		ExecuteActivity(ctx,
			prune.Prune,
			activities.PruneParams{},
		).
		Get(ctx, &pruneResult)

	return &CronResult{RunTime: now}, nil
}
