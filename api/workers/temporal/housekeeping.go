package workers

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Housekeeping(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	w := worker.New(*c, constantstaskqueues.QueueCron, worker.Options{})

	id := fmt.Sprintf("housekeeping_cron_%v", uuid.New())
	options := client.StartWorkflowOptions{
		ID:           id,
		TaskQueue:    constantstaskqueues.QueueCron,
		CronSchedule: "0 0 */2 * *", // every 2 days
	}

	cl := *c

	_, err := cl.ExecuteWorkflow(context.Background(), options, temporalworkflows.BuangHouskeeping)
	if err != nil {
		return nil, err
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
