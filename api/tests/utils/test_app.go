package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
	"skulpture/buang/ports"
	workersinterfaces "skulpture/buang/workers/interfaces"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/negrel/assert"
)

type TestApplicationConfig struct {
	Services              TestApplicationServices
	DurableExecutorConfig DurableExecutorConfiguration
}

type TestApplication struct {
	Url    *string
	http   TestHttpApplication
	config TestApplicationConfig
}

type TestHttpApplication struct {
	chi      *chi.Mux
	Services TestApplicationServices
}

type TestApplicationServices struct {
	Queries              *interfaces.Queries
	GorillaSchemaDecoder *schema.Decoder
	GorillaSchemaEncoder *schema.Encoder
	Docker               client.APIClient
	Workflows            workersinterfaces.Workflows
	Ports                ports.Ports
}

func (a TestApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*TestApplication, error) {
	httpApp := TestHttpApplication{
		chi:      chi,
		Services: a.Services,
	}
	app := TestApplication{
		config: a,
		http:   httpApp,
	}

	assert.True(app.http != TestHttpApplication{}, "app initialized incorrectly")
	assert.True(app.http.Services.Workflows != nil, "workflows must be initialized")

	return &app, nil
}

func (a *TestApplication) Run(ctx context.Context) (func(context.Context), error) {
	ts := httptest.NewUnstartedServer(a.http.chi)
	a.Url = &ts.URL
	ts.Start()
	slog.InfoContext(ctx, fmt.Sprintf("Test HTTP server listening at %s", ts.URL))

	cleanup := func(ctx context.Context) {
		ts.Close()
	}
	return cleanup, nil
}

func (a *TestApplication) GetHttpApplication() *TestHttpApplication {
	return &a.http
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

func (s TestApplicationServices) ToAppApplicationServices() app.ApplicationServices {
	as := app.ApplicationServices(s)

	return as
}
