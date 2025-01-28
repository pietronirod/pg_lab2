package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"service-b/internal/repository"
	"service-b/internal/service"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type CEPResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CEPHandler: Request received")

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "process-cep-handler")
	defer span.End()
	log.Printf("CEPHandler: Received TraceID=%s", span.SpanContext().TraceID().String())

	cep := r.URL.Path[len("/cep/"):]
	log.Printf("CEPHandler: CEP received: %s", cep)

	if len(cep) != 8 {
		log.Printf("CEPHandler: Invalid CEP: %s", cep)
		span.SetStatus(codes.Error, "Invalid CEP length")
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	span.SetAttributes(attribute.String("cep", cep))

	city, err := repository.FetchCityFromCEP(ctx, cep)
	if err != nil {
		log.Printf("CEPHandler: Error fetching city for CEP %s: %v", cep, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error fetching city")
		if errors.Is(err, repository.ErrCEPNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, "error fetching city", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("CEPHandler: City found for CEP %s: %s", cep, city)
	span.SetAttributes(attribute.String("city", city))

	tempC, err := repository.FetchTemperature(ctx, city)
	if err != nil {
		log.Printf("CEPHandler: Error fetching temperature for city %s: %v", city, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error fetching temperature")
		http.Error(w, "error fetching temperature", http.StatusInternalServerError)
		return
	}

	log.Printf("CEPHandler: Temperature for city %s: %f°C", city, tempC)

	tempF := service.CelsiusToFahrenheit(tempC)
	tempK := service.CelsiusToKelvin(tempC)
	log.Printf("CEPHandler: Converted temperatures for city %s: %f°F, %f°K", city, tempF, tempK)

	response := CEPResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("CEPHandler: Error encoding response: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error encoding response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	span.SetStatus(codes.Ok, "Successfully processed CEP")
	log.Println("CEPHandler: Request processed successfully")
}
