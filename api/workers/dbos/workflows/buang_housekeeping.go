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

type BuangHousekeeping app.ApplicationServices

func (h BuangHousekeeping) BuangHousekeeping(ctx dbos.DBOSContext, scheduledTime time.Time) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang housekeeping panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("buang housekeeping panic")
			}
		}
	}()

	s := app.ApplicationServices(h)

	staleDeployments, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetStaleDeploymentsResult, error) {
			getStaleDeployments := activities.GetStaleDeployments(s)

			res, err := getStaleDeployments.GetStaleDeployments(ctx, activities.GetStaleDeploymentsParams{})
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment(s)

			for _, d := range staleDeployments.Deployments {
				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    d.ProjectID,
					DeploymentId: d.DeploymentID,
				})
				if err != nil {
					return nil, err
				}
			}

			return &activities.BuangDeploymentResult{}, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	_, err = dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.PruneResult, error) {
			prune := activities.Prune(s)

			res, err := prune.Prune(ctx)
			if err != nil {
				return nil, err
			}

			return res, nil
		}, dbos.WithStepMaxRetries(3))

	return true, nil
}
