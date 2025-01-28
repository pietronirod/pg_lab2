package main

import (
	"log"
	"net/http"
	"service-a/internal/delivery"
	"service-a/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	cleanup := tracing.InitTracing("service-a")
	defer cleanup()

	http.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(delivery.CEPHandler), "CEPHandler"))

	log.Println("Service A is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
