package workerstemporalclient

import (
	"context"
	"fmt"
	"log/slog"

	"skulpture/buang/app"
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
			panic(fmt.Sprintf("unable to start worker: %v", err.Error()))
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
