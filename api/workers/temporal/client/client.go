package workerstemporalclient

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	enumsenv "skulpture/buang/enums/env"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
	workerstemporal "skulpture/buang/workers/temporal"

	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	temporallog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

type TemporalWorker func(s app.ApplicationServices, c *temporalclient.Client) (worker.Worker, error)

type TemporalConfig struct {
	Services workersshared.WorkflowServices

	HostPort *string
}

type Client struct {
	c temporalclient.Client
}

var _ workersinterfaces.Workflows = (*Client)(nil)

func New(ctx context.Context, cfg TemporalConfig) (workersinterfaces.Workflows, func(ctx context.Context), error) {
	var opts temporalclient.Options

	opts.Logger = temporallog.NewStructuredLogger(slog.Default())
	if cfg.HostPort != nil {
		opts.HostPort = *cfg.HostPort
	} else {
		var err error

		opts, err = envconfig.LoadDefaultClientOptions()
		if err != nil {
			return nil, nil, err
		}
	}

	tc, err := temporalclient.NewLazyClient(opts)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func(ctx context.Context) {
		tc.Close()
	}

	client := &Client{c: tc}

	run := func(w TemporalWorker) {
		_, err := w(cfg.Services.ToAppServices(), &tc)
		if err != nil {
			p := reflect.ValueOf(w).Pointer()
			f := runtime.FuncForPC(p)

			// in tests the scheduled workers start polling when the test is shutting down
			// and it causes things to panic because the worker can't establish a connection
			// breaks not accomodating test code in application code a little but most tests will not
			// fail for no reason if we don't sleep
			if env, _ := enumsenv.Parse(constantsenvs.GO_ENV.Value()); env != enumsenv.Test {
				panic(fmt.Sprintf("unable to start worker %v: %v", f.Name(), err.Error()))
			}

			if env, _ := enumsenv.Parse(constantsenvs.GO_ENV.Value()); env == enumsenv.Test {
				slog.ErrorContext(ctx, fmt.Sprintf("unable to start worker %v: %v", f.Name(), err.Error()))
			}
		}
	}

	ws := []TemporalWorker{
		workerstemporal.DeploymentWorker,
		workerstemporal.BuangWorker,
		workerstemporal.Housekeeping,
	}

	for _, w := range ws {
		go run(w)
	}

	return client, cleanup, nil
}
