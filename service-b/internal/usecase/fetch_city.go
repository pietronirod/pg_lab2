package usecase

import (
	"context"
	"errors"
	"service-b/internal/repository"
)

type FetchCityService interface {
	Fetch(ctx context.Context, cep string) (string, error)
}

type fetchCityService struct{}

func NewFetchCityService() FetchCityService {
	return &fetchCityService{}
}

func (s *fetchCityService) Fetch(ctx context.Context, cep string) (string, error) {
	city, err := repository.FetchCityFromCEP(ctx, cep)
	if err != nil {
		if errors.Is(err, repository.ErrCEPNotFound) {
			return "", repository.ErrCEPNotFound
		}
		return "", err
	}
	return city, nil
}
