package workers

import (
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func BuangWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	w := worker.New(*c, constantstaskqueues.QueueBuang, worker.Options{})

	buangDeployment := activities.BuangDeployment(s)
	activeDeployments := activities.ActiveDeploymentIds(s)

	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(activeDeployments.GetActiveDeploymentIds)
	w.RegisterWorkflow(temporalworkflows.BuangDeployment)
	w.RegisterWorkflow(temporalworkflows.BuangBranch)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
