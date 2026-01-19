package temporalworkflows

import (
	"context"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BuangHousekeeping app.ApplicationServices

type CronResult struct {
	RunTime time.Time
}

func (bh BuangHousekeeping) BuangHousekeeping(ctx workflow.Context) (*CronResult, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang housekeeping panic: %v", r))
		}
	}()

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

	err := workflow.
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

	return &CronResult{RunTime: now}, nil
}
