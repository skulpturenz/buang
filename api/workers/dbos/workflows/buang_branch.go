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

type BuangBranch app.ApplicationServices

type BuangBranchParams struct {
	ProjectId int64
	Branch    string
}

func (bb BuangBranch) BuangBranch(ctx dbos.DBOSContext, p BuangBranchParams) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang branch panic: %v", r))

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("buang deployment panic")
			}
		}
	}()

	s := app.ApplicationServices(bb)

	activeDeploymentIds, err := dbos.RunAsStep(ctx,
		func(ctx context.Context) (*activities.GetActiveDeploymentIdsResult, error) {
			getActiveDeploymentIds := activities.ActiveDeploymentIds(s)
			res, err := getActiveDeploymentIds.GetActiveDeploymentIds(ctx, activities.GetActiveDeploymentIdsParams{
				ProjectId: p.ProjectId,
				Branch:    p.Branch,
			})
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
			errorDeployment := activities.ErrorDeployment(s)

			for _, id := range activeDeploymentIds.DeploymentIds {
				_, err := buangDeployment.BuangDeployment(ctx, activities.BuangDeploymentParams{
					ProjectId:    p.ProjectId,
					DeploymentId: id,
				})
				if err != nil {
					_, errErrorDeployment := errorDeployment.ErrorDeployment(ctx, activities.ErrorDeploymentParams{
						ProjectId:    p.ProjectId,
						DeploymentId: id,
					})
					if errErrorDeployment != nil {
						return nil, errors.Join(err, errErrorDeployment)
					}

					return nil, err
				}
			}

			return &activities.BuangDeploymentResult{}, nil
		}, dbos.WithStepMaxRetries(3))
	if err != nil {
		return false, err
	}

	return true, nil
}
