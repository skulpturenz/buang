package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	enumsenv "skulpture/buang/enums/env"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdklog "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type LoggerConfig struct {
	Enable   bool
	Env      enumsenv.Environment
	Service  string
	LogLevel slog.Leveler
}

func (c LoggerConfig) SetDefault(ctx context.Context, r *chi.Mux) func(context.Context) {
	if !c.Enable {
		// https://github.com/open-telemetry/opentelemetry-go/discussions/2659#discussioncomment-10798740
		otel.SetTracerProvider(
			noop.NewTracerProvider(),
		)

		return func(context.Context) { // noop cleanup
		}
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(),
	)

	if err != nil {
		slog.ErrorContext(ctx, "error", "otel", fmt.Sprintf("failed to create exporter: %s", err.Error()))
		panic(err)
	}
	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", c.Service),
			attribute.String("library.language", "go"),
			attribute.String("deployment.environment", c.Env.String()),
		),
	)
	if err != nil {
		slog.ErrorContext(ctx, "error", "otel", fmt.Sprintf("could not set resources: %s", err.Error()))
		panic(err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			// TODO: reconsider later
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.25))),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	logExporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient(
		// parseable only supports json payloads
		// see: https://www.parseable.com/docs/OpenTelemetry/logs
		otlplogshttp.WithJsonProtocol())))
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithBatcher(logExporter),
		sdklog.WithResource(resources),
	)

	otelLogger := slog.New(
		slogmulti.Fanout(
			otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{
				Level: c.LogLevel,
			}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: c.LogLevel,
			}),
		),
	)
	slog.SetDefault(otelLogger)

	r.Use(otelhttp.NewMiddleware(c.Service))
	r.Use(httplog.RequestLogger(httplog.NewLogger(c.Service, httplog.Options{
		Concise: true,
		Tags: map[string]string{
			"env": c.Env.String(),
		},
	})))

	return func(ctx context.Context) {
		loggerErr := loggerProvider.Shutdown((ctx))
		exporterErr := exporter.Shutdown(ctx)
		joined := errors.Join(loggerErr, exporterErr)

		slog.ErrorContext(ctx, joined.Error())
	}
}
