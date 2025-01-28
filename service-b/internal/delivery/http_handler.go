package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"service-b/internal/repository"
	"service-b/internal/usecase"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// CEPResponse define o formato da resposta do serviço
type CEPResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// CEPHandler gerencia as requisições para buscar cidade e temperatura
type CEPHandler struct {
	fetchCity usecase.FetchCityService
	fetchTemp usecase.FetchTempService
}

// NewCEPHandler cria um novo handler
func NewCEPHandler(fetchCity usecase.FetchCityService, fetchTemp usecase.FetchTempService) *CEPHandler {
	return &CEPHandler{fetchCity: fetchCity, fetchTemp: fetchTemp}
}

// Handle processa a requisição para buscar cidade e temperatura pelo CEP
func (h *CEPHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("CEPHandler: Request received")

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "process-cep-handler")
	defer span.End()
	log.Printf("CEPHandler: Received TraceID=%s", span.SpanContext().TraceID().String())

	cep := r.URL.Path[len("/cep/"):]
	log.Printf("CEPHandler: CEP received: %s", cep)

	// Validação do CEP
	if len(cep) != 8 {
		log.Printf("CEPHandler: Invalid CEP: %s", cep)
		span.SetStatus(codes.Error, "Invalid CEP length")
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}
	span.SetAttributes(attribute.String("cep", cep))

	// Buscar cidade pelo CEP
	city, err := h.fetchCity.Fetch(ctx, cep)
	if err != nil {
		log.Printf("CEPHandler: Error fetching city for CEP %s: %v", cep, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error fetching city")
		if errors.Is(err, repository.ErrCEPNotFound) {
			h.writeErrorResponse(w, http.StatusNotFound, "can not find zipcode")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "error fetching city")
		}
		return
	}

	log.Printf("CEPHandler: City found for CEP %s: %s", cep, city)
	span.SetAttributes(attribute.String("city", city))

	// Buscar temperatura para a cidade
	tempC, err := h.fetchTemp.Fetch(ctx, city)
	if err != nil {
		log.Printf("CEPHandler: Error fetching temperature for city %s: %v", city, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error fetching temperature")
		h.writeErrorResponse(w, http.StatusInternalServerError, "error fetching temperature")
		return
	}

	log.Printf("CEPHandler: Temperature for city %s: %f°C", city, tempC)

	// Converte as temperaturas
	tempF := usecase.CelsiusToFahrenheit(tempC)
	tempK := usecase.CelsiusToKelvin(tempC)
	log.Printf("CEPHandler: Converted temperatures for city %s: %f°F, %f°K", city, tempF, tempK)

	// Monta a resposta
	response := CEPResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	h.writeSuccessResponse(w, response)
	span.SetStatus(codes.Ok, "Successfully processed CEP")
	log.Println("CEPHandler: Request processed successfully")
}

// writeErrorResponse padroniza as respostas de erro
func (h *CEPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeSuccessResponse padroniza as respostas de sucesso
func (h *CEPHandler) writeSuccessResponse(w http.ResponseWriter, response CEPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
