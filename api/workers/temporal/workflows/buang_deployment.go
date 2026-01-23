package temporalworkflows

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"
	"time"

	"github.com/DataDog/gostackparse"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BuangDeployment(ctx workflow.Context, projectId int64, deploymentId int64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang deployment panic: %v", r))

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
				err = errors.New("buang deployment panic")
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

	var buangDeployment *activities.BuangDeployment
	var buangDeploymentResult activities.BuangDeploymentResult

	err = workflow.ExecuteActivity(ctx,
		buangDeployment.BuangDeployment,
		activities.BuangDeploymentParams{
			ProjectId:    projectId,
			DeploymentId: deploymentId,
		}).
		Get(ctx, &buangDeploymentResult)
	if err != nil {
		var errorDeployment *activities.ErrorDeployment
		var errorDeploymentResult activities.ErrorDeploymentResult

		errErrorDeployment := workflow.ExecuteActivity(ctx,
			errorDeployment.ErrorDeployment,
			activities.ErrorDeploymentParams{
				ProjectId:    projectId,
				DeploymentId: deploymentId,
			}).
			Get(ctx, &errorDeploymentResult)
		if errErrorDeployment != nil {
			return errors.Join(err, errErrorDeployment)
		}

		return err
	}

	return nil
}
