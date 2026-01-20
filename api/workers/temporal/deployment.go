package workers

import (
	"context"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func DeploymentWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("deployment panic: %v", r))
		}
	}()

	w := worker.New(*c, constantstaskqueues.QueueDeployment, worker.Options{})

	createDynamicConfigDir := activities.CreateDynamicConfigDir(s)
	cloneDeployment := activities.CloneDeployment(s)
	deployProject := activities.DeployProject(s)
	activeDeploymentIds := activities.ActiveDeploymentIds(s)
	buangDeployment := activities.BuangDeployment(s)
	getDeploymentBranch := activities.GetDeploymentBranch(s)

	w.RegisterActivity(createDynamicConfigDir.CreateDynamicConfigDir)
	w.RegisterActivity(cloneDeployment.CloneDeployment)
	w.RegisterActivity(deployProject.DeployProject)
	w.RegisterActivity(activeDeploymentIds.GetActiveDeploymentIds)
	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(getDeploymentBranch.GetDeploymentBranch)

	w.RegisterWorkflow(temporalworkflows.Deploy)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
