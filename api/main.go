package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"skulpture/buang/app"
	apploggerotel "skulpture/buang/app/logger"
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
	mcphandlers "skulpture/buang/handlers/mcp"
	"skulpture/buang/handlers/projects"
	authn "skulpture/buang/middleware/authn"
	limiter "skulpture/buang/middleware/limiter"
	"skulpture/buang/ports"
	"skulpture/buang/workers"
	workersinterfaces "skulpture/buang/workers/interfaces"
	workersshared "skulpture/buang/workers/shared"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-cz/devslog"
	slogchi "github.com/samber/slog-chi"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	"gitlab.com/greyxor/slogor"

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
	env, err := enumsenv.Parse(constantsenvs.GO_ENV.Value())

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RealIP)

	// TODO: otel might need more work
	otelSloggerCfg := apploggerotel.OtelSloggerConfig{
		Enable:   constantsenvs.ENABLE_TELEMETRY.Value() || env == enumsenv.Production,
		Service:  constantsenvs.OTEL_SERVICE_NAME.Value(),
		Env:      env,
		LogLevel: constantsenvs.LOG_LEVEL.Value(),
	}
	otelSlogger, otelSloggerCleanup, err := otelSloggerCfg.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer otelSloggerCleanup(ctx)

	if v, ok := constantsenvs.BUANG_SENTRY_DSN.Value(); ok {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:           v,
			EnableTracing: true,
			Environment:   env.String(),
			ServerName:    constantsenvs.OTEL_SERVICE_NAME.Value(),
		})
		if err != nil {
			slog.ErrorContext(ctx, "error", "err", err.Error())
			panic(err)
		}
	}

	handlers := []slog.Handler{
		otelSlogger,
	}
	if env == enumsenv.Production {
		handlers = append(handlers,
			slogor.NewHandler(os.Stdout,
				slogor.SetLevel(constantsenvs.LOG_LEVEL.Value()),
				slogor.SetTimeFormat(time.Stamp),
				slogor.ShowSource()))
	}
	if env != enumsenv.Production {
		handlers = append(handlers, devslog.NewHandler(os.Stdout, nil))
	}
	if _, ok := constantsenvs.BUANG_SENTRY_DSN.Value(); ok {
		handler := slogsentry.Option{Level: constantsenvs.LOG_LEVEL.Value()}.
			NewSentryHandler()

		handlers = append(handlers, handler)
	}
	slogger := slog.
		New(slogmulti.Fanout(handlers...)).
		With("environment", env.String()).
		With("release", constantsenvs.BUANG_VERSION)
	slog.SetDefault(slogger)

	r.Use(otelSloggerCfg.Handler)
	r.Use(slogchi.NewWithConfig(slogger, slogchi.Config{
		WithSpanID:         true,
		WithTraceID:        true,
		WithRequestID:      true,
		WithRequestHeader:  constantsenvs.LOG_LEVEL.Value() == slog.LevelDebug,
		WithResponseHeader: constantsenvs.LOG_LEVEL.Value() == slog.LevelDebug,
		WithRequestBody:    constantsenvs.LOG_LEVEL.Value() == slog.LevelDebug,
		DefaultLevel:       constantsenvs.LOG_LEVEL.Value(),
	}))

	r.Use(middleware.Recoverer)

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

	dbType, err := enumsdbtypes.Parse(constantsenvs.DB_TYPE.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	dbCfg := db.DbConfig{
		Type:             dbType,
		ConnectionString: constantsenvs.DB_CONNECTION_STRING.Value(),
	}
	queries, queriesCleanup, err := dbCfg.New(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer queriesCleanup(ctx)

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer docker.Close()

	p := ports.New()
	ws := workersshared.WorkflowServices{
		Queries:              &queries,
		Docker:               docker,
		GorillaSchemaDecoder: schema.NewDecoder(),
		GorillaSchemaEncoder: schema.NewEncoder(),
		Ports:                p,
	}
	workflows, cleanup, err := createWorkflows(ctx, ws)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}
	defer cleanup(ctx)

	s := app.ApplicationServices{
		Queries:              &queries,
		Docker:               docker,
		GorillaSchemaDecoder: schema.NewDecoder(),
		GorillaSchemaEncoder: schema.NewEncoder(),
		Workflows:            workflows,
		Ports:                p,
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

	appConfig := app.ApplicationConfig{
		HttpPort: ":80",
		Services: s,
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

	r.Route("/mcp/api/v1", func(r chi.Router) {
		r.Use(limiter.Handle)
		r.Use(authnConfig.Handle)

		app.GetHttpApplication().AddRouters(r, mcphandlers.Router)
	})

	err = app.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	c := func(ctx context.Context) {
		if v, ok := os.LookupEnv(constantsenvs.INTERNAL_BUANG_BOOTSTRAP_DIR); ok && isExperimentalBootstrapEnabled {
			os.RemoveAll(v)
		}
	}
	defer c(ctx)
}

func createWorkflows(ctx context.Context, services workersshared.WorkflowServices) (workersinterfaces.Workflows, func(context.Context), error) {
	durableExecutor, err := enumsdurableexecutors.Parse(constantsenvs.DURABLE_EXECUTOR.Value())
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	createTemporalWorkflows := func(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
		tc := workers.TemporalConfig{
			Services: services,
		}
		workflows, cleanup, err := tc.New(ctx)
		if err != nil {
			return nil, nil, err
		}

		return workflows, cleanup, nil
	}

	createDbosWorkflows := func(ctx context.Context) (workersinterfaces.Workflows, func(context.Context), error) {
		dc := workers.DbosConfig{
			Services:    services,
			AppName:     constantsenvs.OTEL_SERVICE_NAME.Value(),
			DatabaseURL: constantsenvs.DB_CONNECTION_STRING.Value(),
		}
		if v, ok := constantsenvs.DBOS_CONDUCTOR_API_KEY.Value(); ok {
			dc.ConductorAPIKey = &v
		}
		if v, ok := constantsenvs.DBOS_CONDUCTOR_URL.Value(); ok {
			dc.ConductorURL = &v
		}
		if v, ok := constantsenvs.DBOS_ADMIN_SERVER_PORT.Value(); ok {
			dc.AdminServerPort = &v
		}

		workflows, cleanup, err := dc.New(ctx)
		if err != nil {
			return nil, nil, err
		}

		return workflows, cleanup, nil
	}

	if durableExecutor == enumsdurableexecutors.Temporal {
		return createTemporalWorkflows(ctx)
	}

	if durableExecutor == enumsdurableexecutors.Dbos {
		return createDbosWorkflows(ctx)
	}

	return nil, nil, fmt.Errorf("unsupported durable executor")
}
