package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/db/interfaces"
	"skulpture/buang/ports"
	"skulpture/buang/services"
	workersinterfaces "skulpture/buang/workers/interfaces"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
)

type ApplicationConfig struct {
	HttpPort string
	Services *services.Services[any, any]
}

type Application struct {
	http   HttpApplication
	config ApplicationConfig
}

type HttpApplication struct {
	chi      *chi.Mux
	Services *services.Services[any, any]
}

type ApplicationServices struct {
	*services.Services[any, any]
}

func (a ApplicationServices) GetQueries() (interfaces.Queries, bool) {
	svc, ok := services.Get[interfaces.Queries](a.Services, services.KeyQueries)
	return svc.Unwrap(), ok
}

func (a ApplicationServices) GetGorillaSchemaDecoder() (schema.Decoder, bool) {
	svc, ok := services.Get[schema.Decoder](a.Services, services.KeyGorillaSchemaDecoder)
	return svc.Unwrap(), ok
}

func (a ApplicationServices) GetGorillaSchemaEncoder() (schema.Encoder, bool) {
	svc, ok := services.Get[schema.Encoder](a.Services, services.KeyGorillaSchemaEncoder)
	return svc.Unwrap(), ok
}

func (a ApplicationServices) GetDocker() (*client.Client, bool) {
	svc, ok := services.Get[*client.Client](a.Services, services.KeyDocker)
	return svc.Unwrap(), ok
}

func (a ApplicationServices) GetWorkflows() (workersinterfaces.Workflows, bool) {
	svc, ok := services.Get[workersinterfaces.Workflows](a.Services, services.KeyWorkflows)
	return svc.Unwrap(), ok
}

func (a ApplicationServices) GetPorts() (ports.Ports, bool) {
	svc, ok := services.Get[ports.Ports](a.Services, services.KeyPorts)
	return svc.Unwrap(), ok
}

func (a ApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*Application, error) {
	services := a.Services

	httpApp := HttpApplication{
		chi:      chi,
		Services: services,
	}
	app := Application{
		config: a,
		http:   httpApp,
	}

	assert.True(app.http != HttpApplication{}, "app initialized incorrectly")
	appSvc := ApplicationServices{Services: app.http.Services}
	_, ok := appSvc.GetWorkflows()
	assert.True(ok, "workflows must be initialized")

	return &app, nil
}

func (a *Application) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.config.HttpPort,
		Handler: a.http.chi,
	}
	slog.InfoContext(ctx, "buang running", "version", constantsenvs.BUANG_VERSION)
	slog.InfoContext(ctx, fmt.Sprintf("HTTP server listening on %s", a.config.HttpPort))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { // blocks until server shutdown
		slog.ErrorContext(ctx, err.Error())
	}

	return nil
}

func (a *Application) GetHttpApplication() *HttpApplication {
	return &a.http
}

func (a *Application) AddSingletons(xs ...AppServiceSingleton) {
	for _, x := range xs {
		x.AssertInitialized()
	}
}

type HttpRouter func(s ApplicationServices, r chi.Router)

func (a *HttpApplication) AddRouters(r chi.Router, x ...HttpRouter) {
	for _, y := range x {
		appSvc := ApplicationServices{Services: a.Services}
		y(appSvc, r)
	}
}

type AppServiceSingleton interface {
	AssertInitialized()
}
