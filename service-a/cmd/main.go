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

func main() {
	// Carregar a configuração
	cfg := config.LoadConfig()

	// Inicializar tracing
	shutdown := tracing.InitTracing("service-a")
	defer shutdown()

	// Configurar o propagador de contexto
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Criar instância do handler passando os valores corretamente
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	handler := delivery.NewCEPHandler(cfg.ServiceBURL, httpClient)

	mux := http.NewServeMux()
	mux.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handler.Handle), "cep-handler"))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
