package main

import (
	"log"
	"net/http"
	"service-b/internal/config"
	"service-b/internal/delivery"
	"service-b/internal/repository"
	"service-b/internal/tracing"
	"service-b/internal/usecase"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	// Carregar a configuração
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inicializar tracing
	shutdown := tracing.InitTracing("service-b")
	defer shutdown()

	// Configurar o propagador de contexto
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Criar instâncias dos repositórios
	cityRepo := repository.NewCityRepository()
	tempRepo := repository.NewTemperatureRepository()

	// Criar instâncias dos casos de uso
	fetchCityService := usecase.NewFetchCityService(cityRepo)
	fetchTempService := usecase.NewFetchTempService(tempRepo)

	// Criar instância do handler passando os valores corretamente
	handler := delivery.NewCEPHandler(fetchCityService, fetchTempService)

	mux := http.NewServeMux()
	mux.Handle("/cep/", otelhttp.NewHandler(http.HandlerFunc(handler.Handle), "cep-handler"))

	log.Println("Starting server on :8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
