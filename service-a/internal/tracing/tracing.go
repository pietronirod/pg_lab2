package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracing inicializa o rastreamento com OTEL
func InitTracing(serviceName string) func() {
	ctx := context.Background()

	// Configurar o exportador OTLP
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create OTLP trace exporter: %v", err)
	}

	// Configurar o provedor de rastreamento
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown TracerProvider: %v", err)
		}
	}
}
