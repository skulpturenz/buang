package workers

import (
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/activities"
	"skulpture/buang/workers/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func DeploymentWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	w := worker.New(*c, constantstaskqueues.QueueDeployment, worker.Options{})

	createDynamicConfigDir := activities.CreateDynamicConfigDir(s)
	cloneDeployment := activities.CloneDeployment(s)
	deployProject := activities.DeployProject(s)

	w.RegisterActivity(createDynamicConfigDir)
	w.RegisterActivity(cloneDeployment)
	w.RegisterActivity(deployProject)
	w.RegisterWorkflow(workflows.Deploy)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
