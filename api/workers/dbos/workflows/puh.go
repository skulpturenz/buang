package dbosworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	"skulpture/buang/workers/activities"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type PeriodicUpdateHandler app.ApplicationServices

func (p PeriodicUpdateHandler) PeriodicUpdateHandler(ctx dbos.DBOSContext, scheduledTime time.Time) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("autoupdate buang panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("autoupdate buang panic")
			}
		}
	}()

	s := app.ApplicationServices(p)

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.AutoupdateBuangResult, error) {
			autoupdateBuang := activities.AutoupdateBuang(s)

			res, err := autoupdateBuang.AutoupdateBuang(ctx)
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	return true, nil
}
