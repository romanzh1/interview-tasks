package tracing

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	HTTPTransport = "http"
	GRPCTransport = "grpc"
)

func InitTracer(url string, service string, transport string) (*sdktrace.TracerProvider, error) {
	var client otlptrace.Client
	var err error

	if transport == HTTPTransport {
		client = otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(url),
			otlptracehttp.WithInsecure(),
		)
	} else if transport == GRPCTransport {
		client = otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(url),
			otlptracegrpc.WithInsecure())
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func ShutdownTracer(tp *sdktrace.TracerProvider) {
	if err := tp.Shutdown(context.Background()); err != nil {
		slog.Error("Failed to shutdown tracer", "error", err)
	}
}
