package workers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	constantsfeaturetoggles "skulpture/buang/constants/feature_toggles"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	enumsenv "skulpture/buang/enums/env"
	"skulpture/buang/workers/activities"
	temporalworkflows "skulpture/buang/workers/temporal/workflows"
	"strconv"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Housekeeping(s app.ApplicationServices, c *client.Client) (worker.Worker, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(context.Background(), fmt.Sprintf("housekeeping panic: %v", r))

			log := map[string]any{
				"stack": debug.Stack(),
			}

			p := o11y.CreateDiagnosticLogParams{
				Type: enumsdiagnosticlogtype.Panic,
				Log:  log,
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

	goEnv, _ := os.LookupEnv("GO_ENV")
	goEnvE, _ := enumsenv.Parse(goEnv)
	isExperimentalBootstrapEnabledEnv, _ := os.LookupEnv(constantsfeaturetoggles.EXPERIMENTAL_BOOTSTRAP)
	isExperimentalBootstrapEnabled, _ := strconv.ParseBool(isExperimentalBootstrapEnabledEnv)

	if goEnvE == enumsenv.Production && isExperimentalBootstrapEnabled {
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
