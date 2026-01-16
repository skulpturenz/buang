package workers

import (
	"skulpture/buang/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Worker(s app.ApplicationServices, c *client.Client) worker.Worker {
	w := worker.New(*c, "sdfvs", worker.Options{})

	w.Run(worker.InterruptCh())

	return w
}
