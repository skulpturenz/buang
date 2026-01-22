package temporalworkflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangBranch(ctx workflow.Context, projectId int64, branch string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang branch panic: %v", r))

			log := map[string]any{
				"stack": string(debug.Stack()),
			}

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  log,
			}
			p.Exec(context.Background())

			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("buang branch panic")
			}
		}
	}()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{UhOh.String()},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var activeDeployments *activities.ActiveDeploymentIds
	var getActiveDeploymentsResult activities.GetActiveDeploymentIdsResult

	err = workflow.
		ExecuteActivity(ctx,
			activeDeployments.GetActiveDeploymentIds,
			activities.GetActiveDeploymentIdsParams{
				ProjectId: projectId,
				Branch:    branch,
			}).
		Get(ctx, &getActiveDeploymentsResult)
	if err != nil {
		return err
	}

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	for _, id := range getActiveDeploymentsResult.DeploymentIds {
		err := workflow.
			ExecuteActivity(ctx,
				buangDeployment.BuangDeployment,
				activities.BuangDeploymentParams{
					ProjectId:    projectId,
					DeploymentId: id,
				},
			).
			Get(ctx, &buangDeploymentResult)

		if err != nil {
			var errorDeployment *activities.ErrorDeployment
			var errorDeploymentResult activities.ErrorDeploymentResult

			errErrorDeployment := workflow.ExecuteActivity(ctx,
				errorDeployment.ErrorDeployment,
				activities.ErrorDeploymentParams{
					ProjectId:    projectId,
					DeploymentId: id,
				}).
				Get(ctx, &errorDeploymentResult)
			if errErrorDeployment != nil {
				return errors.Join(err, errErrorDeployment)
			}

			return err
		}
	}

	return nil
}
