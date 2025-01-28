package main

import (
	"log"
	"net/http"
	"service-b/internal/config"
	"service-b/internal/delivery"
	"service-b/internal/repository"
	"service-b/internal/usecase"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	config.LoadConfig()

	// Criar instâncias dos repositórios
	cityRepo := repository.NewCityRepository()
	tempRepo := repository.NewTemperatureRepository()

	// Criar instâncias dos casos de uso
	fetchCityService := usecase.NewFetchCityService(cityRepo)
	fetchTempService := usecase.NewFetchTempService(tempRepo)

	// Criar instância do handler passando os valores corretamente
	handler := delivery.NewCEPHandler(*fetchCityService, *fetchTempService)

	mux := http.NewServeMux()
	mux.Handle("/cep/", otelhttp.NewHandler(http.HandlerFunc(handler.Handle), "cep-handler"))

	log.Println("Service B is running on port 8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
