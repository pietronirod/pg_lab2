package main

import (
	"log"
	"net/http"
	"service-b/internal/delivery"
	"service-b/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	// Configure o propagador para o contexto W3C
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	cleanup := tracing.InitTracing("service-b")
	defer cleanup()

	mux := http.NewServeMux()
	mux.Handle("/cep/", otelhttp.NewHandler(http.HandlerFunc(delivery.CEPHandler), "cep-handler"))

	log.Println("Service B is running on port 8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
