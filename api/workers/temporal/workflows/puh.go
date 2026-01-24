package temporalworkflows

import (
	"skulpture/buang/app"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type PeriodicUpdateHandler app.ApplicationServices

type PeriodicUpdateHandlerResult struct {
	RunTime time.Time
}

func (p PeriodicUpdateHandler) PeriodicUpdateHandler(ctx workflow.Context) (res *PeriodicUpdateHandlerResult, err error) {
	defer workers.RecoverWorkflowPanic(enumsdurableexecutors.Temporal, "PeriodicUpdateHandler", err)

	s := app.ApplicationServices(p)

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	now := workflow.Now(ctx)

	autoupdateBuang := activities.AutoupdateBuang(s)
	var autoupdateBuangResult activities.AutoupdateBuangResult

	err = workflow.
		ExecuteActivity(ctx,
			autoupdateBuang.AutoupdateBuang).
		Get(ctx, &autoupdateBuangResult)
	if err != nil {
		return nil, err
	}

	return &PeriodicUpdateHandlerResult{RunTime: now}, nil
}
