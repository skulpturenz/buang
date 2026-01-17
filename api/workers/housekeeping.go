package workers

import (
	"context"
	"fmt"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/workflows"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Housekeeping(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	id := fmt.Sprintf("housekeeping_cron_%v", uuid.New())
	options := client.StartWorkflowOptions{
		ID:           id,
		TaskQueue:    constantstaskqueues.QueueCron,
		CronSchedule: "0 0 */2 * *", // every 2 days
	}

	cl := *c

	_, err := cl.ExecuteWorkflow(context.Background(), options, workflows.BuangHouskeeping)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
