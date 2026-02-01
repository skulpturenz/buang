package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/db/interfaces"
	"skulpture/buang/ports"
	workersinterfaces "skulpture/buang/workers/interfaces"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
)

type ApplicationConfig struct {
	HttpPort string
	Services ApplicationServices
}

type Application struct {
	http   HttpApplication
	config ApplicationConfig
}

type HttpApplication struct {
	chi      *chi.Mux
	Services ApplicationServices
}

type ApplicationServices struct {
	Queries              *interfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               client.APIClient
	Workflows            workersinterfaces.Workflows
	Ports                ports.Ports
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
	assert.True(app.http.Services.Workflows != nil, "workflows must be initialized")

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
		y(a.Services, r)
	}
}

type AppServiceSingleton interface {
	AssertInitialized()
}
