package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	apppkg "skulpture/buang/app"
	"skulpture/buang/services"

	"github.com/go-chi/chi/v5"
	"github.com/negrel/assert"
)

type TestApplicationConfig struct {
	Services              *services.Services[any, any]
	DurableExecutorConfig DurableExecutorConfiguration
}

type TestApplication struct {
	Url    *string
	http   TestHttpApplication
	config TestApplicationConfig
}

type TestHttpApplication struct {
	chi      *chi.Mux
	Services *services.Services[any, any]
}

type TestApplicationServices struct {
	*services.Services[any, any]
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
	appSvc := apppkg.ApplicationServices{Services: app.http.Services}
	_, ok := appSvc.GetWorkflows()
	assert.True(ok, "workflows must be initialized")

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

func (a *TestApplication) AddSingletons(xs ...apppkg.AppServiceSingleton) {
	for _, x := range xs {
		x.AssertInitialized()
	}
}

type HttpRouter func(s apppkg.ApplicationServices, r chi.Router)

func (a *TestHttpApplication) AddRouters(r chi.Router, x ...apppkg.HttpRouter) {
	for _, y := range x {
		y(apppkg.ApplicationServices{Services: a.Services}, r)
	}
}

func (s TestApplicationServices) ToAppApplicationServices() apppkg.ApplicationServices {
	return apppkg.ApplicationServices{Services: s.Services}
}
