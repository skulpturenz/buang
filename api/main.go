package main

import (
	"context"
	"log/slog"
	"os"
	"skulpture/buang/app"
	"skulpture/buang/components/o11y"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/db"
	_ "skulpture/buang/docs"
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	enumsenv "skulpture/buang/enums/env"
	deploymentlogs "skulpture/buang/handlers/deployment_logs"
	"skulpture/buang/handlers/deployments"
	"skulpture/buang/handlers/diagnostics"
	"skulpture/buang/handlers/projects"
	authn "skulpture/buang/middleware/authn"
	limiter "skulpture/buang/middleware/limiter"
	"skulpture/buang/workers/dbos"
	workers "skulpture/buang/workers/temporal"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title						Buang API
// @description				Deploy preview environments with ease
// @license					MIT
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-Key
func main() {
	ctx := context.Background()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	env, err := enumsenv.Parse(constantsenvs.GO_ENV.Value())
	limiterConfig := limiter.LimiterConfig{
		Env:      env,
		Tokens:   500,
		Interval: time.Minute,
	}
	limiter, err := limiterConfig.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	logger := app.LoggerConfig{
		Enable:   constantsenvs.ENABLE_TELEMETRY.Value() || env == enumsenv.Production,
		Service:  constantsenvs.OTEL_SERVICE_NAME.Value(),
		Env:      env,
		LogLevel: constantsenvs.LOG_LEVEL.Value(),
	}
	cleanup := logger.SetDefault(ctx, r)
	defer cleanup(ctx)

	dbType, err := enumsdbtypes.Parse(constantsenvs.DB_TYPE.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	dbCfg := db.DbConfig{
		Type:             dbType,
		ConnectionString: constantsenvs.DB_CONNECTION_STRING.Value(),
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

	s := app.ApplicationServices{
		Queries:              &queries,
		Docker:               docker,
		GorillaSchemaDecoder: schema.NewDecoder(),
		GorillaSchemaEncoder: schema.NewEncoder(),
	}

	isExperimentalBootstrapEnabled, ok := constantsenvs.EXPERIMENTAL_BOOTSTRAP.Value()
	if constantsenvs.GO_ENV.Value() == enumsenv.Production.String() && ok && isExperimentalBootstrapEnabled { // TODO: test
		err := bootstrap(ctx, &s)
		if err != nil {
			panic(err)
		}

		return
	}

	authnConfig := authn.AuthnMiddlewareConfig{
		ApiKey: constantsenvs.API_KEY.Value(),
	}

	durableExecutor, err := enumsdurableexecutors.Parse(constantsenvs.DURABLE_EXECUTOR.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	appConfig := app.ApplicationConfig{
		HttpPort:        ":80",
		Services:        s,
		DurableExecutor: durableExecutor,
	}
	app, err := appConfig.New(ctx, r)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	app.AddSingletons(o11y.NewS(&queries))

	r.Mount("/docs", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {

		r.Use(limiter.Handle)
		r.Use(authnConfig.Handle)

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

	c := func(ctx context.Context) {
		cleanup(ctx)

		if v, ok := os.LookupEnv(constantsenvs.INTERNAL_BUANG_BOOTSTRAP_DIR); ok && isExperimentalBootstrapEnabled {
			os.RemoveAll(v)
		}
	}

	defer c(ctx)
}
