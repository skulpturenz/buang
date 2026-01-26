package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"skulpture/buang/components/o11y"
	"skulpture/buang/db"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	deploymentlogs "skulpture/buang/handlers/deployment_logs"
	"skulpture/buang/handlers/deployments"
	"skulpture/buang/handlers/diagnostics"
	"skulpture/buang/handlers/projects"
	"skulpture/buang/workers"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
)

func CreateApp(ctx context.Context, cfg TestApplicationConfig) (*TestApplication, func(context.Context)) {
	r := chi.NewRouter()

	dbCfg := db.DbConfig{ // db choice constrained by the type of durable executor used
		Type:             cfg.DurableExecutorConfig.DbType,
		ConnectionString: cfg.DurableExecutorConfig.DbConnectionString,
	}
	queries, dbCleanup, err := dbCfg.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer dbCleanup(ctx)

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer docker.Close()

	app, err := cfg.New(ctx, r)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	app.AddSingletons(o11y.NewS(&queries))

	r.Route("/api/v1", func(r chi.Router) {
		app.GetHttpApplication().AddRouters(r, projects.Router,
			deployments.Router,
			deploymentlogs.Router,
			diagnostics.Router)
	})

	cleanup, err := app.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	return app, cleanup
}

func CreateWorkflows(ctx context.Context, cfg DurableExecutorConfiguration, services workersshared.WorkflowServices) (workersinterfaces.Workflows, func(context.Context), error) {
	createTemporalWorkflows := func(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
		c := *cfg.DurableExecutorContainer
		port, err := c.MappedPort(ctx, "7233/tcp")
		if err != nil {
			panic(err)
		}

		hostPort := fmt.Sprintf("%v:%v", "0.0.0.0", port.Port())
		tc := workers.TemporalConfig{
			Services: services,
			HostPort: &hostPort,
		}

		workflows, cleanup, err := tc.New(ctx)

		if err != nil {
			return nil, cleanup, err
		}

		return workflows, cleanup, nil
	}

	createDbosWorkflows := func(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
		dc := workers.DbosConfig{
			Services:    services,
			AppName:     "test",
			DatabaseURL: cfg.DbConnectionString,
		}

		workflows, cleanup, err := dc.New(ctx)
		if err != nil {
			return nil, cleanup, err
		}

		return workflows, cleanup, nil
	}

	if cfg.DurableExecutor == enumsdurableexecutors.Temporal {
		return createTemporalWorkflows(ctx)
	}

	if cfg.DurableExecutor == enumsdurableexecutors.Dbos {
		return createDbosWorkflows(ctx)
	}

	return nil, nil, fmt.Errorf("unsupported durable executor")
}
