package dbosworkflows

import (
	"context"
	"skulpture/buang/app"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"skulpture/buang/workers/activities"
	workersshared "skulpture/buang/workers/shared"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type PeriodicUpdateHandler app.ApplicationServices

func (p PeriodicUpdateHandler) PeriodicUpdateHandler(ctx dbos.DBOSContext, scheduledTime time.Time) (res bool, err error) {
	defer workersshared.RecoverWorkflowPanic(enumsdurableexecutors.Dbos, "PeriodicUpdateHandler", err)

	s := app.ApplicationServices(p)

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.AutoupdateBuangResult, error) {
			autoupdateBuang := activities.AutoupdateBuang(s)

			res, err := autoupdateBuang.AutoupdateBuang(ctx)
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(10))
	if err != nil {
		return false, err
	}

	return true, nil
}
