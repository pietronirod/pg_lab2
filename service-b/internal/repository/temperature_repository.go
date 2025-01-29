package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"service-b/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TemperatureRepository interface {
	FetchTemperature(ctx context.Context, city string) (float64, error)
}

type temperatureRepository struct{}

// NewTemperatureRepository cria um novo repositório TemperatureRepository
func NewTemperatureRepository() TemperatureRepository {
	return &temperatureRepository{}
}

// FetchTemperature busca a temperatura de uma cidade
func (r *temperatureRepository) FetchTemperature(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "fetch-temperature")
	defer span.End()

	apiKey := config.AppConfig.WeatherAPIKey
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s?key=%s&q=%s", config.AppConfig.WeatherAPIURL, apiKey, encodedCity)

	log.Printf("FetchTemperature: Fetching temperature for city: %s", city)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create request")
		log.Printf("Error creating request: %v", err)
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to fetch temperature")
		log.Printf("Error fetching temperature: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, fmt.Sprintf("Non-OK HTTP status: %s", resp.Status))
		log.Printf("Non-OK HTTP status: %s", resp.Status)
		return 0, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		span.SetStatus(codes.Error, "Failed to decode response")
		log.Printf("Error decoding response: %v", err)
		return 0, err
	}

	span.SetAttributes(attribute.String("city", city), attribute.Float64("temperature", result.Current.TempC))
	span.SetStatus(codes.Ok, "Successfully fetched temperature")
	return result.Current.TempC, nil
}
