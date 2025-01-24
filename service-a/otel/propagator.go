package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}

func InitPropagator() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
