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
)

type TemperatureRepository interface {
	FetchTemperature(ctx context.Context, city string) (float64, error)
}

type temperatureRepository struct{}

func NewTemperatureRepository() TemperatureRepository {
	return &temperatureRepository{}
}

func (r *temperatureRepository) FetchTemperature(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "fetch-temperature")
	defer span.End()

	apiKey := config.AppConfig.WeatherAPIKey
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s?key=%s&q=%s", config.AppConfig.WeatherAPIURL, apiKey, encodedCity)

	log.Printf("FetchTemperature: Fetching temperature for city: %s", city)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("FetchTemperature: Error making request to WeatherAPI: %v", err)
		span.RecordError(err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("FetchTemperature: WeatherAPI returned non-200 status: %d", resp.StatusCode)
		return 0, fmt.Errorf("failed to fetch temperature for city: %s", city)
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("FetchTemperature: Error decoding response: %v", err)
		span.RecordError(err)
		return 0, err
	}

	span.SetAttributes(attribute.Float64("fetch.temperature_c", result.Current.TempC))

	return result.Current.TempC, nil
}
