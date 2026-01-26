package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"os/signal"
	"reflect"
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
	workers  []app.TemporalWorker
}

type TestDbosApplication struct {
	ctx       dbos.DBOSContext
	Services  TestApplicationServices
	workflows []app.DbosWorkflow[any, any]
}

type TestApplicationServices struct {
	Queries              *interfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               *client.Client
	durableExecutor      enumsdurableexecutors.DurableExecutor
	temporal             *temporalclient.Client
	dbos                 *dbos.DBOSContext
}

// NOTE: MOSTLY FOLLOWS THE REAL THING
// EXCEPT FOR HOW WE CREATE THE HTTP SERVER
// AND SOME TYPE MAPPING
func (a TestApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*TestApplication, error) {
	services := a.Services
	services.durableExecutor = a.DurableExecutorConfig.DurableExecutor

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
		services.temporal = &hc
		httpApp.Services.temporal = &hc

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
		services.dbos = &dbosContext
		httpApp.Services.dbos = &dbosContext

		app.http = httpApp
		app.dbos = &TestDbosApplication{
			ctx:      dbosContext,
			Services: services,
		}
	}

	assert.True(app.http != TestHttpApplication{}, "app initialized incorrectly")
	assert.True(app.config.DurableExecutorConfig.DurableExecutor != enumsdurableexecutors.Dbos || app.http.Services.dbos != nil, "dbos is injected to be used by handlers")
	assert.True(app.config.DurableExecutorConfig.DurableExecutor != enumsdurableexecutors.Temporal || app.http.Services.temporal != nil, "temporal is injected to be used by handlers")
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

		run := func(w app.TemporalWorker) {
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

func (a *TestHttpApplication) AddRouters(r chi.Router, x ...app.HttpRouter) {
	for _, y := range x {
		y(a.Services.ToAppApplicationServices(), r)
	}
}

func (a *TestTemporalApplication) AddWorkers(ws ...app.TemporalWorker) {
	a.workers = ws
}

func (a *TestDbosApplication) AddWorkflows(ws ...app.DbosWorkflow[any, any]) {
	a.workflows = ws
}

func (s TestApplicationServices) ToAppApplicationServices() app.ApplicationServices {
	as := app.ApplicationServices{
		Queries:              s.Queries,
		GorillaSchemaDecoder: s.GorillaSchemaDecoder,
		GorillaSchemaEncoder: s.GorillaSchemaEncoder,
		Docker:               s.Docker,
	}

	// this is a hack just for tests
	// durable executors should not be initialized from the outside, they are set and configured internally
	// don't want to have a setter because i think we might as well make it public in that case
	// and move instantiation of durable executors to the top level
	// but with temporal workers since we're not deploying them independently we want to control how they're run
	// and also so that app instantiation is simpler, without workflows app does not function
	// want to create like:
	//  - create `ApplicationServices`
	//  - create `New` `Application` from those services
	//  - attach workflows and handlers to that
	// proper handling: similar to how dbs works. `Workflows` interface and each durable executor would have to
	// implement it. with workflows for each type of durable executor attach there instead of the app itself as it is now
	setPrivateField(&as, "durableExecutor", s.durableExecutor)
	setPrivateField(&as, "temporal", s.temporal)
	setPrivateField(&as, "dbos", s.dbos)

	return as
}

func setPrivateField[T app.ApplicationServices](v *T, fieldName string, fieldValue any) *T {
	val := reflect.ValueOf(v).Elem()

	// based on: https://medium.com/@darshan.na185/modifying-private-variables-of-a-struct-in-go-using-unsafe-and-reflect-5447b3019a80
	member := val.FieldByName(fieldName)
	// basically get the pointer to the field on the struct
	// and update the value that it points to
	// so from zero values for the fields to the new value
	// - the fields pointer on the struct stays the same, but the value the pointer points to has changed
	// ... so we make a new copy of `fieldValue` and update the contents of the existing pointer?
	ptrToMember := reflect.NewAt(member.Type(), member.Addr().UnsafePointer()).Elem()
	if ptrToMember.IsValid() && ptrToMember.CanSet() {
		ptrToMember.Set(reflect.ValueOf(fieldValue))
	}

	return v
}
