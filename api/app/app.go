package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/db/interfaces"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"strconv"
	"sync"
	"syscall"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	temporallog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

type ApplicationConfig struct {
	HttpPort        string
	Services        ApplicationServices
	DurableExecutor enumsdurableexecutors.DurableExecutor
}

type Application struct {
	http     HttpApplication
	temporal *TemporalApplication
	config   ApplicationConfig
	dbos     *DbosApplication
}

type HttpApplication struct {
	chi      *chi.Mux
	Services ApplicationServices
}

type TemporalApplication struct {
	client   *temporalclient.Client
	Services ApplicationServices
	workers  []TemporalWorker
}

type DbosApplication struct {
	ctx       dbos.DBOSContext
	Services  ApplicationServices
	workflows []DbosWorkflow[any, any]
}

type ApplicationServices struct {
	Queries              *interfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               *client.Client
	durableExecutor      enumsdurableexecutors.DurableExecutor
	temporal             *temporalclient.Client
	dbos                 dbos.DBOSContext
}

func (a ApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*Application, error) {
	services := a.Services
	services.durableExecutor = a.DurableExecutor

	httpApp := HttpApplication{
		chi:      chi,
		Services: services,
	}
	app := Application{
		config: a,
	}

	// so that attaching workers does not throw
	initialTemporal := &TemporalApplication{}
	initialDbos := &DbosApplication{}
	app.temporal = initialTemporal
	app.dbos = initialDbos

	if a.DurableExecutor == enumsdurableexecutors.Temporal {
		opts := envconfig.MustLoadDefaultClientOptions()
		opts.Logger = temporallog.NewStructuredLogger(slog.Default())

		tc, err := temporalclient.NewLazyClient(opts) // used by the workers
		if err != nil {
			return nil, err
		}

		hc, err := temporalclient.NewLazyClient(opts) // used by the routes
		if err != nil {
			return nil, err
		}
		services.temporal = &hc
		httpApp.Services.temporal = &hc

		app.http = httpApp
		app.temporal = &TemporalApplication{
			Services: services,
			client:   &tc,
		}
	} else {
		cfg := dbos.Config{
			AppName:     os.Getenv("OTEL_SERVICE_NAME"),
			DatabaseURL: os.Getenv("DB_CONNECTION_STRING"),
		}
		if os.Getenv("DBOS_CONDUCTOR_API_KEY") != "" {
			cfg.ConductorAPIKey = os.Getenv("DBOS_CONDUCTOR_API_KEY")
		}
		if os.Getenv("DBOS_CONDUCTOR_URL") != "" {
			cfg.ConductorURL = os.Getenv("DBOS_CONDUCTOR_URL")
		}
		if os.Getenv("DBOS_ADMIN_SERVER_PORT") != "" {
			p := os.Getenv("DBOS_ADMIN_SERVER_PORT")
			port, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}

			cfg.AdminServer = true
			cfg.AdminServerPort = port
		}

		dbosContext, err := dbos.NewDBOSContext(context.Background(), cfg)
		if err != nil {
			return nil, err
		}
		services.dbos = dbosContext
		httpApp.Services.dbos = dbosContext

		app.http = httpApp
		app.dbos = &DbosApplication{
			ctx:      dbosContext,
			Services: services,
		}
	}

	assert.True(app.http != HttpApplication{}, "app initialized incorrectly")
	assert.True(app.config.DurableExecutor != enumsdurableexecutors.Dbos || app.http.Services.dbos != nil, "dbos is injected to be used by handlers")
	assert.True(app.config.DurableExecutor != enumsdurableexecutors.Temporal || app.http.Services.temporal != nil, "temporal is injected to be used by handlers")
	assert.True(app.temporal != initialTemporal || app.dbos != initialDbos, "must use one durable executor")

	assert.True(app.config.DurableExecutor != enumsdurableexecutors.Temporal || app.temporal.client != nil,
		"must have a client if temporal application")

	return &app, nil
}

func (a *Application) Run(ctx context.Context) (func(ctx context.Context), error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	startHttpServer := func(ctx context.Context, wg *sync.WaitGroup) {
		server := &http.Server{
			Addr:    a.config.HttpPort,
			Handler: a.http.chi,
		}

		go func() {
			defer wg.Done()
			slog.InfoContext(ctx, fmt.Sprintf("HTTP server listening on %s", a.config.HttpPort))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { // blocks until server shutdown
				slog.ErrorContext(ctx, err.Error())
			}

		}()

		_ = <-ctx.Done()
		slog.InfoContext(ctx, "shutting down HTTP server")
		if err := server.Shutdown(ctx); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("server shutdown failed: %v", err.Error()))
		}
	}

	runTemporalWorkers := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		run := func(w TemporalWorker) {
			wg.Add(1)
			defer wg.Done()

			_, err := w(a.temporal.Services, a.temporal.client) // blocks until workers stop
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("unable to start worker: %v", err.Error()))

				cancel()
			}
		}

		for _, w := range a.temporal.workers {
			go run(w)
		}

		slog.InfoContext(ctx, "started temporal workers")
	}

	addDbosWorkflows := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for _, w := range a.dbos.workflows {
			w(a.dbos.Services, a.dbos.ctx)
		}

		err := dbos.Launch(a.dbos.ctx)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("unable to start dbos: %v", err.Error()))

			cancel()
		}
	}

	wg.Add(1)
	go startHttpServer(ctx, &wg)

	if a.config.DurableExecutor == enumsdurableexecutors.Temporal {
		wg.Add(1)
		go runTemporalWorkers(ctx, &wg)

	} else {
		wg.Add(1)
		go addDbosWorkflows(ctx, &wg)
	}

	go func() {
		_ = <-sigs
		cancel()
	}()

	slog.InfoContext(ctx, "buang running", "version", constantsenvs.BUANG_VERSION)

	wg.Wait()

	cleanup := func(ctx context.Context) {
		c := *a.temporal.client

		c.Close()
	}

	return cleanup, nil
}

func (a *Application) GetHttpApplication() *HttpApplication {
	return &a.http
}

func (a *Application) GetTemporalApplication() *TemporalApplication {
	return a.temporal
}

func (a *Application) GetDbosApplication() *DbosApplication {
	return a.dbos
}

func (a *Application) AddSingletons(xs ...AppServiceSingleton) {
	for _, x := range xs {
		x.AssertInitialized()
	}
}

type HttpRouter func(s ApplicationServices, r chi.Router)

func (a *HttpApplication) AddRouters(r chi.Router, x ...HttpRouter) {
	for _, y := range x {
		y(a.Services, r)
	}
}

type TemporalWorker func(s ApplicationServices, c *temporalclient.Client) (worker.Worker, error)

func (a *TemporalApplication) AddWorkers(ws ...TemporalWorker) {
	a.workers = ws
}

type DbosWorkflow[P any, R any] func(s ApplicationServices, c dbos.DBOSContext)

func (a *DbosApplication) AddWorkflows(ws ...DbosWorkflow[any, any]) {
	a.workflows = ws
}

func (s *ApplicationServices) GetDurableExecutor() (any, enumsdurableexecutors.DurableExecutor) {
	if s.durableExecutor == enumsdurableexecutors.Temporal {
		return *s.temporal, s.durableExecutor
	}

	return s.dbos, s.durableExecutor
}

type TemporalClient = temporalclient.Client

type DbosContext = dbos.DBOSContext

type AppServiceSingleton interface {
	AssertInitialized()
}
