package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}

func InitPropagator(serviceName, collectorURL string) func() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return func() {}
}
