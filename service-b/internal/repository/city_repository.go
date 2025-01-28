package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"service-b/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var ErrCEPNotFound = errors.New("CEP not found")

type CityRepository interface {
	FetchCityFromCEP(ctx context.Context, cep string) (string, error)
}

type cityRepository struct{}

func NewCityRepository() CityRepository {
	return &cityRepository{}
}

func (r *cityRepository) FetchCityFromCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "fetch-city-from-cep")
	defer span.End()

	url := fmt.Sprintf("%s%s/json/", config.AppConfig.ViaCEPAPIURL, cep)
	log.Printf("FetchCityFromCEP: Fetching city for CEP %s from URL: %s", cep, url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("FetchCityFromCEP: Error making request to ViaCEP: %v", err)
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("FetchCityFromCEP: Non-200 status code returned")
		return "", ErrCEPNotFound
	}

	var result struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("FetchCityFromCEP: Error decoding response: %v", err)
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(attribute.String("fetch.city", result.Localidade))

	if result.Localidade == "" {
		log.Println("FetchCityFromCEP: City not found in response")
		return "", ErrCEPNotFound
	}

	return result.Localidade, nil
}
