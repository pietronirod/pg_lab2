package usecase_test

import (
	"context"
	"service-b/internal/repository"
	"service-b/internal/usecase"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock para o repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FetchCityFromCEP(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}

func TestFetchCityService_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := usecase.NewFetchCityService(mockRepo)

	cep := "01001000"
	expectedCity := "São Paulo"

	// Configuração do mock
	mockRepo.On("FetchCityFromCEP", mock.Anything, cep).Return(expectedCity, nil)

	// Execução do teste
	city, err := service.Fetch(context.Background(), cep)

	// Validação
	require.NoError(t, err)
	require.Equal(t, expectedCity, city)

	mockRepo.AssertExpectations(t)
}

func TestFetchCityService_CEPNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := usecase.NewFetchCityService(mockRepo)

	cep := "99999999"

	// Configuração do mock para erro de CEP não encontrado
	mockRepo.On("FetchCityFromCEP", mock.Anything, cep).Return("", repository.ErrCEPNotFound)

	// Execução do teste
	city, err := service.Fetch(context.Background(), cep)

	// Validação
	require.ErrorIs(t, err, repository.ErrCEPNotFound)
	require.Empty(t, city)

	mockRepo.AssertExpectations(t)
}
