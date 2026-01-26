package apploggerotel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	enumsenv "skulpture/buang/enums/env"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdklog "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type OtelSloggerConfig struct {
	Enable   bool
	Env      enumsenv.Environment
	Service  string
	LogLevel slog.Leveler
}

func (c OtelSloggerConfig) New(ctx context.Context) (slog.Handler, func(context.Context), error) {
	if !c.Enable {
		// https://github.com/open-telemetry/opentelemetry-go/discussions/2659#discussioncomment-10798740
		otel.SetTracerProvider(
			noop.NewTracerProvider(),
		)

		handler := slog.NewTextHandler(io.Discard, nil)

		cleanup := func(context.Context) {} // noop

		return handler, cleanup, nil
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

func (c OtelSloggerConfig) Handler(next http.Handler) http.Handler {
	return otelhttp.NewMiddleware(c.Service)(next)
}
