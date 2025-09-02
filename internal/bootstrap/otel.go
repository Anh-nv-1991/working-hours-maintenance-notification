package bootstrap

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Tracing struct {
	Shutdown func(context.Context) error
}

func InitTracing(ctx context.Context, serviceName string) *Tracing {
	exp, err := otlptracehttp.New(ctx) // đọc endpoint từ ENV: OTEL_EXPORTER_OTLP_ENDPOINT
	if err != nil {
		log.Fatalf("otlp exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		log.Fatalf("resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return &Tracing{Shutdown: tp.Shutdown}
}
