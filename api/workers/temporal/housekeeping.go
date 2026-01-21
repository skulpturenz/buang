package workers

import (
	"context"
	"fmt"
	"log/slog"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Housekeeping(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("housekeeping panic: %v", r))
		}
	}()

	w := worker.New(*c, constantstaskqueues.QueueCron, worker.Options{})

	getStaleDeployments := activities.GetStaleDeployments(s)
	buangDeployment := activities.BuangDeployment(s)
	prune := activities.Prune(s)
	w.RegisterActivity(getStaleDeployments.GetStaleDeployments)
	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(prune.Prune)

	buangHousekeeping := temporalworkflows.BuangHousekeeping(s)
	w.RegisterWorkflow(buangHousekeeping.BuangHousekeeping)

	id := fmt.Sprintf("housekeeping_cron_%v", uuid.New())
	options := client.StartWorkflowOptions{
		ID:           id,
		TaskQueue:    constantstaskqueues.QueueCron,
		CronSchedule: "0 0 */2 * *", // every 2 days
	}

	cl := *c
	_, err := cl.ExecuteWorkflow(context.Background(), options, buangHousekeeping.BuangHousekeeping)
	if err != nil {
		return nil, err
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
