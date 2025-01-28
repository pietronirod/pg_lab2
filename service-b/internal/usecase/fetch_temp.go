package usecase

import (
	"context"
	"service-b/internal/repository"
)

type FetchTempService struct {
	repo repository.TemperatureRepository
}

func NewFetchTempService(repo repository.TemperatureRepository) *FetchTempService {
	return &FetchTempService{repo: repo}
}

func (s *FetchTempService) Fetch(ctx context.Context, city string) (float64, error) {
	return s.repo.FetchTemperature(ctx, city)
}
