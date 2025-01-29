package usecase

import (
	"context"
	"log"
	"service-b/internal/repository"
)

// FetchCityService define a interface para buscar a cidade correspondente a um CEP
type FetchCityService interface {
	Fetch(ctx context.Context, cep string) (string, error)
}

type fetchCityService struct {
	repo repository.CityRepository
}

// NewFetchCityService cria um novo serviço FetchCityService
func NewFetchCityService(repo repository.CityRepository) FetchCityService {
	return &fetchCityService{repo: repo}
}

// Fetch busca a cidade correspondente a um CEP
func (s *fetchCityService) Fetch(ctx context.Context, cep string) (string, error) {
	city, err := s.repo.FetchCityFromCEP(ctx, cep)
	if err != nil {
		log.Printf("Error fetching city for CEP %s: %v", cep, err)
		return "", err
	}
	return city, nil
}
