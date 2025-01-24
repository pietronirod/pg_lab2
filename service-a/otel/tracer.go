package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer(serviceName, collectorURL string) func() {
	exporter, err := zipkin.New(collectorURL)
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", serviceName),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		_ = tracerProvider.Shutdown(context.Background())
	}
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("example").Start(ctx, name)
}
