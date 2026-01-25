package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"os/signal"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/db/interfaces"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	"sync"
	"syscall"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
	temporalclient "go.temporal.io/sdk/client"
	temporallog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

type TestApplicationConfig struct {
	Services              TestApplicationServices
	DurableExecutorConfig DurableExecutorConfiguration
}

type TestApplication struct {
	Url      *string
	http     TestHttpApplication
	temporal *TestTemporalApplication
	config   TestApplicationConfig
	dbos     *TestDbosApplication
}

type TestHttpApplication struct {
	chi      *chi.Mux
	Services TestApplicationServices
}

type TestTemporalApplication struct {
	client   *temporalclient.Client
	Services TestApplicationServices
	workers  []TemporalWorker
}

type TestDbosApplication struct {
	ctx       dbos.DBOSContext
	Services  TestApplicationServices
	workflows []DbosWorkflow[any, any]
}

type TestApplicationServices struct {
	Queries              *interfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               *client.Client
	DurableExecutor      enumsdurableexecutors.DurableExecutor
	Temporal             *temporalclient.Client
	Dbos                 dbos.DBOSContext
}

// NOTE: MOSTLY FOLLOWS THE REAL THING
// EXCEPT FOR HOW WE CREATE THE HTTP SERVER
// AND SOME TYPE MAPPING
func (a TestApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*TestApplication, error) {
	services := a.Services
	services.DurableExecutor = a.DurableExecutorConfig.DurableExecutor

	httpApp := TestHttpApplication{
		chi:      chi,
		Services: services,
	}
	app := TestApplication{
		config: a,
	}

	// so that attaching workers does not throw
	initialTemporal := &TestTemporalApplication{}
	initialDbos := &TestDbosApplication{}
	app.temporal = initialTemporal
	app.dbos = initialDbos

	if a.DurableExecutorConfig.DurableExecutor == enumsdurableexecutors.Temporal {
		c := *a.DurableExecutorConfig.DurableExecutorContainer
		port, err := c.MappedPort(ctx, "7233/tcp")
		if err != nil {
			return nil, err
		}

		opts := temporalclient.Options{
			HostPort: fmt.Sprintf("%v:%v", "0.0.0.0", port.Port()),
		}
		opts.Logger = temporallog.NewStructuredLogger(slog.Default())

		tc, err := temporalclient.NewLazyClient(opts) // used by the workers
		if err != nil {
			return nil, err
		}

		hc, err := temporalclient.NewLazyClient(opts) // used by the routes
		if err != nil {
			return nil, err
		}
		services.Temporal = &hc
		httpApp.Services.Temporal = &hc

		app.http = httpApp
		app.temporal = &TestTemporalApplication{
			Services: services,
			client:   &tc,
		}
	} else {
		cfg := dbos.Config{
			AppName:     "test",
			DatabaseURL: a.DurableExecutorConfig.DbConnectionString,
		}

		dbosContext, err := dbos.NewDBOSContext(context.Background(), cfg)
		if err != nil {
			return nil, err
		}
		services.Dbos = dbosContext
		httpApp.Services.Dbos = dbosContext

		app.http = httpApp
		app.dbos = &TestDbosApplication{
			ctx:      dbosContext,
			Services: services,
		}
	}

	assert.True(app.http != TestHttpApplication{}, "app initialized incorrectly")
	assert.True(app.config.DurableExecutorConfig.DurableExecutor != enumsdurableexecutors.Dbos || app.http.Services.Dbos != nil, "dbos is injected to be used by handlers")
	assert.True(app.config.DurableExecutorConfig.DurableExecutor != enumsdurableexecutors.Temporal || app.http.Services.Temporal != nil, "temporal is injected to be used by handlers")
	assert.True(app.temporal != initialTemporal || app.dbos != initialDbos, "must use one durable executor")
	assert.True(app.config.DurableExecutorConfig.DurableExecutor != enumsdurableexecutors.Temporal || app.temporal.client != nil,
		"must have a client if temporal application")

	return &app, nil
}

func (a *TestApplication) Run(ctx context.Context) (func(ctx context.Context), error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	ts := httptest.NewUnstartedServer(a.http.chi)
	a.Url = &ts.URL
	ts.Start()
	slog.InfoContext(ctx, fmt.Sprintf("Test HTTP server listening at %s", ts.URL))

	runTemporalWorkers := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		run := func(w TemporalWorker) {
			wg.Add(1)
			defer wg.Done()

			_, err := w(a.temporal.Services.ToAppApplicationServices(), a.temporal.client) // blocks until workers stop
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
			w(a.temporal.Services.ToAppApplicationServices(), a.dbos.ctx)
		}

		err := dbos.Launch(a.dbos.ctx)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("unable to start dbos: %v", err.Error()))

			cancel()
		}
	}

	if a.config.DurableExecutorConfig.DurableExecutor == enumsdurableexecutors.Temporal {
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
		if a.temporal.client != nil {
			c := *a.temporal.client

			c.Close()
		}

		ts.Close()
	}

	return cleanup, nil
}

func (a *TestApplication) GetHttpApplication() *TestHttpApplication {
	return &a.http
}

func (a *TestApplication) GetTemporalApplication() *TestTemporalApplication {
	return a.temporal
}

func (a *TestApplication) GetDbosApplication() *TestDbosApplication {
	return a.dbos
}

func (a *TestApplication) AddSingletons(xs ...app.AppServiceSingleton) {
	for _, x := range xs {
		x.AssertInitialized()
	}
}

type HttpRouter func(s app.ApplicationServices, r chi.Router)

func (a *TestHttpApplication) AddRouters(r chi.Router, x ...HttpRouter) {
	for _, y := range x {
		y(a.Services.ToAppApplicationServices(), r)
	}
}

type TemporalWorker func(s app.ApplicationServices, c *temporalclient.Client) (worker.Worker, error)

func (a *TestTemporalApplication) AddWorkers(ws ...TemporalWorker) {
	a.workers = ws
}

type DbosWorkflow[P any, R any] func(s app.ApplicationServices, c dbos.DBOSContext)

func (a *TestDbosApplication) AddWorkflows(ws ...DbosWorkflow[any, any]) {
	a.workflows = ws
}

func (s *TestApplicationServices) GetDurableExecutor() (any, enumsdurableexecutors.DurableExecutor) {
	if s.DurableExecutor == enumsdurableexecutors.Temporal {
		return *s.Temporal, s.DurableExecutor
	}

	return s.Dbos, s.DurableExecutor
}

func (s TestApplicationServices) ToAppApplicationServices() app.ApplicationServices {
	return app.ApplicationServices{
		Queries:              s.Queries,
		GorillaSchemaDecoder: s.GorillaSchemaDecoder,
		GorillaSchemaEncoder: s.GorillaSchemaEncoder,
		Docker:               s.Docker,
		Temporal:             s.Temporal,
		DurableExecutor:      s.DurableExecutor,
		Dbos:                 s.Dbos,
	}
}
