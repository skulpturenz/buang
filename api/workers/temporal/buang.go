package workers

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/DataDog/gostackparse"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func BuangWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("buang panic: %v", r))

			stack := debug.Stack()
			goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

			log := map[string]any{
				"stack": goroutines,
			}

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  log,
			}
			p.Exec(context.Background())
		}
	}()

	w := worker.New(*c, constantstaskqueues.QueueBuang, worker.Options{})

	buangDeployment := activities.BuangDeployment(s)
	activeDeployments := activities.ActiveDeploymentIds(s)
	errorDeployment := activities.ErrorDeployment(s)

	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(activeDeployments.GetActiveDeploymentIds)
	w.RegisterActivity(errorDeployment.ErrorDeployment)

	w.RegisterWorkflow(temporalworkflows.BuangDeployment)
	w.RegisterWorkflow(temporalworkflows.BuangBranch)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
