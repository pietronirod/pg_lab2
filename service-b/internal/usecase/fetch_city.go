package usecase

import (
	"context"
	"service-b/internal/repository"
)

type FetchCityService struct {
	repo repository.CityRepository
}

func NewFetchCityService(repo repository.CityRepository) *FetchCityService {
	return &FetchCityService{repo: repo}
}

func (s *FetchCityService) Fetch(ctx context.Context, cep string) (string, error) {
	return s.repo.FetchCityFromCEP(ctx, cep)
}
