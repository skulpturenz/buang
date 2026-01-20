package dbosworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	"skulpture/buang/workers/activities"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

type BuangDeployment app.ApplicationServices

type BuangDeploymentParams struct {
	ProjectId    int64
	DeploymentId int64
}

func (bd BuangDeployment) BuangDeployment(ctx dbos.DBOSContext, p BuangDeploymentParams) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang deployment panic: %v", r))
		}
	}()

	s := app.ApplicationServices(bd)

	_, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.BuangDeploymentResult, error) {
			buangDeployment := activities.BuangDeployment(s)

			_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
				ProjectId:    p.ProjectId,
				DeploymentId: p.DeploymentId,
			})
			if err != nil {
				return nil, err
			}

			return &activities.BuangDeploymentResult{}, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		_, errErrorDeployment := dbos.RunAsStep(ctx,
			func(ctx context.Context) (*activities.ErrorDeploymentResult, error) {
				errorDeployment := activities.ErrorDeployment(bd)
				res, err := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: p.DeploymentId,
				})
				if err != nil {
					return nil, err
				}

				return res, err
			}, dbos.WithStepMaxRetries(3))
		if errErrorDeployment != nil {
			return false, errors.Join(err, errErrorDeployment)
		}

		return false, err
	}

	return true, nil
}
