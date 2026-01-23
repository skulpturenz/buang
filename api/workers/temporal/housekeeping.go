package workers

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	constantsenvs "skulpture/buang/constants/envs"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	enumsenv "skulpture/buang/enums/env"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"

	"github.com/DataDog/gostackparse"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Housekeeping(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("housekeeping panic: %v", r))

			stack := debug.Stack()
			goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  goroutines,
			}
			p.Exec(context.Background())
		}
	}()

	w := worker.New(*c, constantstaskqueues.QueueCron, worker.Options{})

	getStaleDeployments := activities.GetStaleDeployments(s)
	buangDeployment := activities.BuangDeployment(s)
	prune := activities.Prune(s)
	autoupdateBuang := activities.AutoupdateBuang(s)
	w.RegisterActivity(getStaleDeployments.GetStaleDeployments)
	w.RegisterActivity(buangDeployment.BuangDeployment)
	w.RegisterActivity(prune.Prune)
	w.RegisterActivity(autoupdateBuang.AutoupdateBuang)

	buangHousekeeping := temporalworkflows.BuangHousekeeping(s)
	periodicUpdateHandler := temporalworkflows.PeriodicUpdateHandler(s)
	w.RegisterWorkflow(buangHousekeeping.BuangHousekeeping)
	w.RegisterWorkflow(periodicUpdateHandler.PeriodicUpdateHandler)

	id := fmt.Sprintf("housekeeping_cron_prune_%v", uuid.New())
	options := client.StartWorkflowOptions{
		ID:           id,
		TaskQueue:    constantstaskqueues.QueueCron,
		CronSchedule: constantsenvs.HOUSEKEEPING_PRUNE_DEPLOYMENTS.Value(),
	}

	cl := *c

	_, err := cl.ExecuteWorkflow(context.Background(), options, buangHousekeeping.BuangHousekeeping)
	if err != nil {
		return nil, err
	}

	goEnv, _ := enumsenv.Parse(constantsenvs.GO_ENV.Value())
	isBootstrapEnabled, _ := constantsenvs.EXPERIMENTAL_BOOTSTRAP.Value()

	if goEnv == enumsenv.Production && isBootstrapEnabled {
		id := fmt.Sprintf("housekeeping_cron_puh_%v", uuid.New())
		options := client.StartWorkflowOptions{
			ID:           id,
			TaskQueue:    constantstaskqueues.QueueCron,
			CronSchedule: constantsenvs.HOUSEKEEPING_PRUNE_DEPLOYMENTS.Value(),
		}

		_, err = cl.ExecuteWorkflow(context.Background(), options, periodicUpdateHandler.PeriodicUpdateHandler)
		if err != nil {
			return nil, err
		}
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		return nil, err
	}

	return w, nil
}
