package usecase

import (
	"context"
	"service-b/internal/repository"
)

type FetchTempService interface {
	Fetch(ctx context.Context, city string) (float64, error)
}

type fetchTempService struct{}

func NewFetchTempService() FetchTempService {
	return &fetchTempService{}
}

func (s *fetchTempService) Fetch(ctx context.Context, city string) (float64, error) {
	temp, err := repository.FetchTemperature(ctx, city)
	if err != nil {
		return 0, err
	}
	return temp, nil
}
