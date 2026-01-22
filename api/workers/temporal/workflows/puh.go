package temporalworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
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
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("periodic update handler panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("periodic update handler panic")
			}
		}
	}()

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
