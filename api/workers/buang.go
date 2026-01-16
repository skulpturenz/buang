package workers

import (
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/activities"
	"skulpture/buang/workers/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func BuangWorker(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	w := worker.New(*c, constantstaskqueues.QueueBuang, worker.Options{})

	buangDeployment := activities.BuangDeployment(s)

	w.RegisterActivity(buangDeployment)
	w.RegisterWorkflow(workflows.Buang)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
