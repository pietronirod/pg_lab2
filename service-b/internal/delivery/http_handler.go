package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"service-b/internal/repository"
	"service-b/internal/usecase"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

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
		if err == repository.ErrCEPNotFound {
			log.Printf("CEPHandler: CEP not found: %s", cep)
			span.SetStatus(codes.Error, "CEP not found")
			h.writeErrorResponse(w, http.StatusNotFound, "can not find zipcode")
		} else {
			log.Printf("CEPHandler: Error fetching city for CEP %s: %v", cep, err)
			span.SetStatus(codes.Error, "Error fetching city")
			h.writeErrorResponse(w, http.StatusInternalServerError, "error fetching city")
		}
		return
	}
	span.SetAttributes(attribute.String("city", city))

	// Buscar temperatura pela cidade
	tempC, err := h.fetchTemp.Fetch(ctx, city)
	if err != nil {
		log.Printf("CEPHandler: Error fetching temperature for city %s: %v", city, err)
		span.SetStatus(codes.Error, "Error fetching temperature")
		h.writeErrorResponse(w, http.StatusInternalServerError, "error fetching temperature")
		return
	}
	span.SetAttributes(attribute.Float64("temperature_celsius", tempC))

	// Converter temperaturas
	tempF := usecase.CelsiusToFahrenheit(tempC)
	tempK := usecase.CelsiusToKelvin(tempC)

	// Responder com a cidade e as temperaturas
	response := map[string]interface{}{
		"city":   city,
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}
	h.writeJSONResponse(w, http.StatusOK, response)
}

func (h *CEPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + message + `"}`))
}

func (h *CEPHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
