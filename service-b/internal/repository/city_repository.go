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
	"go.opentelemetry.io/otel/codes"
)

var ErrCEPNotFound = errors.New("CEP not found")

type CityRepository interface {
	FetchCityFromCEP(ctx context.Context, cep string) (string, error)
}

type cityRepository struct{}

// NewCityRepository cria um novo repositório CityRepository
func NewCityRepository() CityRepository {
	return &cityRepository{}
}

// FetchCityFromCEP busca a cidade correspondente a um CEP
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
		span.SetStatus(codes.Error, "Error making request to ViaCEP")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("non-OK HTTP status: %s", resp.Status)
		log.Printf("FetchCityFromCEP: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Non-OK HTTP status")
		return "", err
	}

	var result struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("FetchCityFromCEP: Error decoding response: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error decoding response")
		return "", err
	}

	if result.Localidade == "" {
		log.Printf("FetchCityFromCEP: CEP %s not found", cep)
		span.SetStatus(codes.Error, "CEP not found")
		return "", ErrCEPNotFound
	}

	span.SetAttributes(attribute.String("cep", cep), attribute.String("city", result.Localidade))
	span.SetStatus(codes.Ok, "Successfully fetched city")
	return result.Localidade, nil
}
