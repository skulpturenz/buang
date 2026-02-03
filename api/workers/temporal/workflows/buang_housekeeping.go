package temporalworkflows

import (
	"context"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers/activities"
	workersshared "skulpture/buang/workers/shared"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BuangHousekeeping app.ApplicationServices

type BuangHousekeepingResult struct {
	RunTime time.Time
}

func (bh BuangHousekeeping) BuangHousekeeping(ctx workflow.Context) (res *BuangHousekeepingResult, err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "Deploy", err)

	s := app.ApplicationServices(bh)

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	now := workflow.Now(ctx)

	getStaleDeployments := activities.GetStaleDeployments(s)
	var getStaleDeploymentsResult activities.GetStaleDeploymentsResult

	err = workflow.
		ExecuteActivity(ctx,
			getStaleDeployments.GetStaleDeployments).
		Get(ctx, &getStaleDeploymentsResult)
	if err != nil {
		return nil, err
	}

	buangDeployment := activities.BuangDeployment(s)
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

	prune := activities.Prune(s)
	var pruneResult activities.PruneResult

	workflow.
		ExecuteActivity(ctx,
			prune.Prune,
			activities.PruneParams{},
		).
		Get(ctx, &pruneResult)

	_ = o11y.PruneStaleDiagnosticLogs(context.Background())

	return &BuangHousekeepingResult{RunTime: now}, nil
}
