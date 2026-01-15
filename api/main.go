package main

import (
	"context"
	"log/slog"
	"net/http"
	"skulpture/buang/enums"
	"time"

	"github.com/dogmatiq/ferrite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

var (
	LOG_LEVEL = ferrite.EnumAs[slog.Level]("LOG_LEVEL", "Log level").
			WithMembers(slog.LevelDebug, slog.LevelError, slog.LevelInfo, slog.LevelWarn).
			WithDefault(slog.LevelInfo).
			Required()
	GO_ENV = ferrite.
		Enum("GO_ENV", "Golang environment").
		WithMembers(enums.Production.String(), enums.Development.String(), enums.Test.String()).
		WithDefault(enums.Development.String()).
		Required()
)

func init() {
	ferrite.Init()
}

func main() {
	ctx := context.Background()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	var store limiter.Store
	if GO_ENV.Value() != enums.Production.String() {
		noopStore, err := noopstore.New()
		if err != nil {
			slog.ErrorContext(ctx, "error", "init", err.Error())
			panic(err)
		}

		store = noopStore
	} else {
		memoryStore, err := memorystore.New(&memorystore.Config{
			Tokens:   5,
			Interval: time.Minute,
		})
		if err != nil {
			slog.ErrorContext(ctx, "error", "init", err.Error())
			panic(err)
		}

		store = memoryStore
	}

	limiter, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc("X-Forwarded-For"))
	if err != nil {
		slog.ErrorContext(ctx, "error", "err", err.Error())
		panic(err)
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(limiter.Handle)

	})

	http.ListenAndServe(":80", r)
}
