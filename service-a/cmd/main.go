package main

import (
	"log"
	"net/http"
	"service-a/internal/config"
	"service-a/internal/delivery"
	"service-a/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	// Configuração do propagador OpenTelemetry
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	// Carregar configurações com Viper
	cfg := config.LoadConfig()

	// Inicializar o tracing distribuído
	cleanup := tracing.InitTracing("service-a")
	defer cleanup()

	// Configurar o handler com a URL do Service B
	handler := delivery.NewCEPHandler(cfg.ServiceBURL)

	// Rota HTTP
	mux := http.NewServeMux()
	mux.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handler.Handle), "CEPHandler"))

	log.Println("Service A is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
