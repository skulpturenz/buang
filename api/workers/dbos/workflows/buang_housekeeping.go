package dbosworkflows

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"
	"time"

	"github.com/DataDog/gostackparse"
	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangHousekeeping app.ApplicationServices

func (h BuangHousekeeping) BuangHousekeeping(ctx dbos.DBOSContext, scheduledTime time.Time) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang housekeeping panic: %v", r))

			stack := debug.Stack()
			goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  goroutines,
			}
			p.Exec(context.Background())

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
	if err != nil {
		return false, err
	}

	_ = o11y.PruneStaleDiagnosticLogs(context.Background())

	return true, nil
}
