package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"service-b/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var ErrCEPNotFound = errors.New("CEP not found")

// FetchCityFromCEP consulta a API ViaCEP para obter a cidade a partir do CEP
func FetchCityFromCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "fetch-city-from-cep")
	defer span.End()

	// Obter URL base do ViaCEP do Viper
	cfg := config.LoadConfig()
	viaCEPURL := cfg.ViaCEPAPIURL
	if viaCEPURL == "" {
		log.Fatal("VIACEP_API_URL is not set in configuration")
	}

	url := fmt.Sprintf("%s%s/json/", viaCEPURL, cep)
	log.Printf("FetchCityFromCEP: Fetching city for CEP %s from URL: %s", cep, url)
	span.SetAttributes(attribute.String("http.url", url), attribute.String("http.method", "GET"))

	// Fazer a requisição HTTP
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("FetchCityFromCEP: Error making request to ViaCEP: %v", err)
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("FetchCityFromCEP: Response status code: %d", resp.StatusCode)
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		log.Println("FetchCityFromCEP: Non-200 status code returned")
		return "", ErrCEPNotFound
	}

	// Decodificar o JSON retornado
	var result struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("FetchCityFromCEP: Error decoding response: %v", err)
		span.RecordError(err)
		return "", err
	}

	log.Printf("FetchCityFromCEP: City found: %s", result.Localidade)
	span.SetAttributes(attribute.String("fetch.city", result.Localidade))

	if result.Localidade == "" {
		log.Println("FetchCityFromCEP: City not found in response")
		return "", ErrCEPNotFound
	}

	return result.Localidade, nil
}

// FetchTemperature consulta a WeatherAPI para obter a temperatura de uma cidade
func FetchTemperature(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "fetch-temperature")
	defer span.End()

	// Obter chave da WeatherAPI e URL base do Viper
	cfg := config.LoadConfig()
	apiKey := cfg.WeatherAPIKey
	weatherAPIURL := cfg.WeatherAPIURL
	if apiKey == "" {
		log.Fatal("WEATHERAPI_KEY is not set in configuration")
	}
	if weatherAPIURL == "" {
		log.Fatal("WEATHERAPI_URL is not set in configuration")
	}

	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s?key=%s&q=%s", weatherAPIURL, apiKey, encodedCity)

	log.Printf("FetchTemperature: Fetching temperature for city: %s", city)
	span.SetAttributes(attribute.String("http.url", url), attribute.String("http.method", "GET"))
	log.Printf("FetchTemperature: Using WeatherAPI URL: %s", url)

	// Fazer a requisição HTTP
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("FetchTemperature: Error making request to WeatherAPI: %v", err)
		span.RecordError(err)
		return 0, err
	}
	defer resp.Body.Close()

	log.Printf("FetchTemperature: Received response with status code: %d", resp.StatusCode)
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		log.Printf("FetchTemperature: WeatherAPI returned non-200 status: %d", resp.StatusCode)
		span.RecordError(fmt.Errorf("failed to fetch temperature for city: %s", city))
		return 0, fmt.Errorf("failed to fetch temperature for city: %s", city)
	}

	// Decodificar o JSON retornado
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

	log.Printf("FetchTemperature: Temperature for city %s: %.2f°C", city, result.Current.TempC)
	span.SetAttributes(attribute.Float64("fetch.temperature_c", result.Current.TempC))

	return result.Current.TempC, nil
}
