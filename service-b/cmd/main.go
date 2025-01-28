package main

import (
	"log"
	"net/http"
	"service-b/internal/config"
	"service-b/internal/delivery"
	"service-b/internal/tracing"
	"service-b/internal/usecase"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	// Configuração do propagador OpenTelemetry
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Carregar configurações
	if config.LoadConfig() == nil {
		log.Fatal("Failed to load configuration")
	}
}

func main() {
	// Inicializar tracing distribuído
	cleanup := tracing.InitTracing("service-b")
	defer cleanup()

	// Criar instâncias dos casos de uso
	fetchCityService := usecase.NewFetchCityService()
	fetchTempService := usecase.NewFetchTempService()

	// Criar handler HTTP e injetar dependências
	cepHandler := delivery.NewCEPHandler(fetchCityService, fetchTempService)

	// Configurar rotas HTTP
	mux := http.NewServeMux()
	mux.Handle("/cep/", otelhttp.NewHandler(http.HandlerFunc(cepHandler.Handle), "cep-handler"))

	log.Println("Service B is running on port 8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
