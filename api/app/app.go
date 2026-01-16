package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"skulpture/buang/db/interfaces"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	temporallog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

type ApplicationConfig struct {
	HttpPort string
	Services ApplicationServices
}

type Application struct {
	http     HttpApplication
	temporal TemporalApplication
	config   ApplicationConfig
}

type HttpApplication struct {
	chi      *chi.Mux
	Services ApplicationServices
}

type TemporalApplication struct {
	client   *client.Client
	Services ApplicationServices
	workers  []TemporalWorker
}

type ApplicationServices struct {
	Queries       *interfaces.Queries
	SchemaDecoder *schema.Decoder
	SchemaEncoder *schema.Encoder
}

func (a ApplicationConfig) New(ctx context.Context, chi *chi.Mux) (*Application, error) {
	opts := envconfig.MustLoadDefaultClientOptions()
	opts.Logger = temporallog.NewStructuredLogger(slog.Default())

	c, err := client.NewLazyClient(opts)
	if err != nil {
		return nil, err
	}

	services := a.Services

	services.SchemaDecoder = schema.NewDecoder()
	services.SchemaEncoder = schema.NewEncoder()

	return &Application{
		http: HttpApplication{
			chi:      chi,
			Services: services,
		},
		temporal: TemporalApplication{
			Services: services,
			client:   &c,
		},
		config: a,
	}, nil
}

func (a *Application) Run(ctx context.Context) (func(ctx context.Context), error) {
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
			log.Fatalf("server shutdown failed: %v", err)
		}
	}

	runTemporalWorkers := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		run := func(w TemporalWorker) {
			wg.Add(1)
			defer wg.Done()

			w(a.temporal.Services, a.temporal.client) // blocks until workers stop
		}

		for _, w := range a.temporal.workers {
			go run(w)
		}

		slog.InfoContext(ctx, "started temporal workers")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	wg.Add(1)
	go startHttpServer(ctx, &wg)

	wg.Add(1)
	go runTemporalWorkers(ctx, &wg)

	go func() {
		_ = <-sigs
		cancel()
	}()

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
	return &a.temporal
}

type HttpRouter func(s ApplicationServices, r chi.Router)

func (a *HttpApplication) AddRouter(r chi.Router, x HttpRouter) {
	x(a.Services, r)
}

type TemporalWorker func(s ApplicationServices, c *client.Client) worker.Worker

func (a *TemporalApplication) AddWorkers(ws ...TemporalWorker) {
	a.workers = ws
}
