package usecase

import (
	"context"
	"log"
	"service-b/internal/repository"
)

// FetchTempService define a interface para buscar a temperatura de uma cidade
type FetchTempService interface {
	Fetch(ctx context.Context, city string) (float64, error)
}

type fetchTempService struct {
	repo repository.TemperatureRepository
}

// NewFetchTempService cria um novo serviço FetchTempService
func NewFetchTempService(repo repository.TemperatureRepository) FetchTempService {
	return &fetchTempService{repo: repo}
}

// Fetch busca a temperatura de uma cidade
func (s *fetchTempService) Fetch(ctx context.Context, city string) (float64, error) {
	temp, err := s.repo.FetchTemperature(ctx, city)
	if err != nil {
		log.Printf("Error fetching temperature for city %s: %v", city, err)
		return 0, err
	}
	return temp, nil
}
