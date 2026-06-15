package testutils

import (
	"context"
	"fmt"
	"log/slog"
	apppkg "skulpture/buang/app"
	"skulpture/buang/components/o11y"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	deploymentlogs "skulpture/buang/handlers/deployment_logs"
	"skulpture/buang/handlers/deployments"
	"skulpture/buang/handlers/diagnostics"
	"skulpture/buang/handlers/projects"
	"skulpture/buang/workers"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"

	"github.com/go-chi/chi/v5"
)

func CreateApp(ctx context.Context, cfg TestApplicationConfig) (*TestApplication, func(context.Context)) {
	r := chi.NewRouter()

	app, err := cfg.New(ctx, r)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	appSvc := apppkg.ApplicationServices{Services: cfg.Services}
	queries, ok := appSvc.GetQueries()
	if !ok {
		panic("queries service not found")
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
