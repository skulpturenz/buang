package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	enumsenv "skulpture/buang/enums/env"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdklog "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
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

func (c LoggerConfig) New(ctx context.Context) (slog.Handler, func(context.Context), error) {
	if !c.Enable {
		// https://github.com/open-telemetry/opentelemetry-go/discussions/2659#discussioncomment-10798740
		otel.SetTracerProvider(
			noop.NewTracerProvider(),
		)

		return nil, func(context.Context) {}, nil
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create exporter: %w", err)
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
		return nil, nil, fmt.Errorf("could not set resources: %w", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			// TODO: reconsider later
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.25))),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	logExporter, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient(
		// parseable only supports json payloads
		// see: https://www.parseable.com/docs/OpenTelemetry/logs
		otlplogshttp.WithJsonProtocol())))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithBatcher(logExporter),
		sdklog.WithResource(resources),
	)

	handler := otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{
		Level: c.LogLevel,
	})

	cleanup := func(ctx context.Context) {
		loggerErr := loggerProvider.Shutdown((ctx))
		exporterErr := exporter.Shutdown(ctx)
		joined := errors.Join(loggerErr, exporterErr)

		if joined != nil {
			slog.ErrorContext(ctx, "error", "otel", joined.Error())
		}
	}

	return handler, cleanup, nil
}

func (c LoggerConfig) AttachMiddleware(r chi.Router) {
	r.Use(otelhttp.NewMiddleware(c.Service))
	r.Use(httplog.RequestLogger(httplog.NewLogger(c.Service, httplog.Options{
		Concise: true,
		Tags: map[string]string{
			"env": c.Env.String(),
		},
	})))
}
