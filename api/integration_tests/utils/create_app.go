package testutils

import (
	"context"
	"log/slog"
	"skulpture/buang/components/o11y"
	"skulpture/buang/db"
	deploymentlogs "skulpture/buang/handlers/deployment_logs"
	"skulpture/buang/handlers/deployments"
	"skulpture/buang/handlers/diagnostics"
	"skulpture/buang/handlers/projects"
	"skulpture/buang/workers/dbos"
	workers "skulpture/buang/workers/temporal"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
)

func CreateApp(ctx context.Context, cfg TestApplicationConfig) (*TestApplication, func(context.Context)) {
	r := chi.NewRouter()

	dbCfg := db.DbConfig{ // db choice constrained by the type of durable executor used
		Type:             cfg.DurableExecutorConfig.DbType,
		ConnectionString: cfg.DurableExecutorConfig.DbConnectionString,
	}
	queries, cleanup, err := dbCfg.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer cleanup(ctx)

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

	app.GetTemporalApplication().AddWorkers(workers.DeploymentWorker,
		workers.BuangWorker,
		workers.Housekeeping,
	)

	app.GetDbosApplication().AddWorkflows(dbos.Deployment,
		dbos.Buang,
		dbos.Housekeping)

	cleanup, err = app.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	return app, cleanup
}
