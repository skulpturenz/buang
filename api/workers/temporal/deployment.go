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

func DeploymentWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("deployment panic: %v", r))

			stack := debug.Stack()
			goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  goroutines,
			}
			p.Exec(context.Background())
		}
	}()

	w := worker.New(*c, constantstaskqueues.QueueDeployment, worker.Options{})

	createDynamicConfigDir := activities.CreateDynamicConfigDir(s)
	cloneDeployment := activities.CloneDeployment(s)
	deployProject := activities.DeployProject(s)
	activeDeploymentIds := activities.ActiveDeploymentIds(s)
	buangDeployment := activities.BuangDeployment(s)
	getDeploymentBranch := activities.GetDeploymentBranch(s)
	errorDeployment := activities.ErrorDeployment(s)

	w.RegisterActivity(createDynamicConfigDir.CreateDynamicConfigDir)
	w.RegisterActivity(cloneDeployment.CloneDeployment)
	w.RegisterActivity(deployProject.DeployProject)
	w.RegisterActivity(activeDeploymentIds.GetActiveDeploymentIds)
	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(getDeploymentBranch.GetDeploymentBranch)
	w.RegisterActivity(errorDeployment.ErrorDeployment)

	w.RegisterWorkflow(temporalworkflows.Deploy)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
