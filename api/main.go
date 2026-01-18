package main

import (
	"context"
	"log/slog"
	"skulpture/buang/app"
	"skulpture/buang/db"
	_ "skulpture/buang/docs"
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	enumsenv "skulpture/buang/enums/env"
	"skulpture/buang/handlers/projects"
	authn "skulpture/buang/middleware/authn"
	limiter "skulpture/buang/middleware/limiter"
	"skulpture/buang/workers/dbos"
	workers "skulpture/buang/workers/temporal"
	"time"

	"github.com/dogmatiq/ferrite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	GO_ENV = ferrite.
		Enum("GO_ENV", "Golang environment").
		WithMembers(enumsenv.Production.String(), enumsenv.Development.String(), enumsenv.Test.String()).
		WithDefault(enumsenv.Development.String()).
		Required()
	API_KEY = ferrite.
		String("API_KEY", "API key").
		WithSensitiveContent().
		WithDefault("supersecureapikey").
		Required()
	LOG_LEVEL = ferrite.EnumAs[slog.Level]("LOG_LEVEL", "Log level").
			WithMembers(slog.LevelDebug, slog.LevelError, slog.LevelInfo, slog.LevelWarn).
			WithDefault(slog.LevelInfo).
			Required()
	ENABLE_TELEMETRY = ferrite.
				Bool("ENABLE_TELEMETRY", "Enable telemetry").
				WithDefault(false).
				Required()
	OTEL_SERVICE_NAME = ferrite.
				String("OTEL_SERVICE_NAME", "OpenTelemetry service name").
				WithDefault("skulpture-buang").
				Required()
	OTEL_EXPORTER_OTLP_ENDPOINT = ferrite.
					String("OTEL_EXPORTER_OTLP_ENDPOINT", "OpenTelemetry exporter endpoint").
					WithDefault("").
					Optional()
	DB_TYPE = ferrite.
		Enum("DB_TYPE", "Database type").
		WithMembers(enumsdbtypes.Pg.String(), enumsdbtypes.Sqlite.String()).
		WithDefault(enumsdbtypes.Sqlite.String()).
		Required()
	DB_CONNECTION_STRING = ferrite.
				String("DB_CONNECTION_STRING", "Database connection string").
				WithSensitiveContent().
				WithDefault("file:test.db?_foreign_keys=true&mode=memory").
				Required()
	TEMPORAL_ADDRESS = ferrite.
				String("TEMPORAL_API_KEY", "Temporal API key").
				WithSensitiveContent().
				Optional()
	TEMPORAL_NAMESPACE = ferrite.
				String("TEMPORAL_NAMESPACE", "Temporal namespace").
				Optional()
	TEMPORAL_API_KEY = ferrite.
				String("TEMPORAL_ADDRESS", "Temporal address").
				Optional()
	DURABLE_EXECUTOR = ferrite.
				Enum("DURABLE_EXECUTOR", "Durable executor").
				WithMembers(enumsdurableexecutors.Temporal.String(), enumsdurableexecutors.Dbos.String()).
				WithDefault(enumsdurableexecutors.Temporal.String()).
				Required()
)

func init() {
	ferrite.Init()
}

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

	env, err := enumsenv.Parse(GO_ENV.Value())
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
		Enable:   ENABLE_TELEMETRY.Value() || env == enumsenv.Production,
		Service:  OTEL_SERVICE_NAME.Value(),
		Env:      env,
		LogLevel: LOG_LEVEL.Value(),
	}
	cleanup := logger.SetDefault(ctx, r)
	defer cleanup(ctx)

	dbType, err := enumsdbtypes.Parse(DB_TYPE.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	dbCfg := db.DbConfig{
		Type:             dbType,
		ConnectionString: DB_CONNECTION_STRING.Value(),
	}
	queries, cleanup, err := dbCfg.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer cleanup(ctx)

	authnConfig := authn.AuthnMiddlewareConfig{
		ApiKey: API_KEY.Value(),
	}

	durableExecutor, err := enumsdurableexecutors.Parse(DURABLE_EXECUTOR.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	appConfig := app.ApplicationConfig{
		HttpPort: ":80",
		Services: app.ApplicationServices{
			Queries: &queries,
		},
		DurableExecutor: durableExecutor,
	}
	app, err := appConfig.New(ctx, r)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	r.Mount("/docs", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {

		r.Use(limiter.Handle)
		r.Use(authnConfig.Handle)

		app.GetHttpApplication().AddRouters(r, projects.Router)
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

	defer cleanup(ctx)
}
