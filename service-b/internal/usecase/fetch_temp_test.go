package usecase_test

import (
	"context"
	"errors"
	"service-b/internal/usecase"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock para o repository
type MockTempRepository struct {
	mock.Mock
}

func (m *MockTempRepository) FetchTemperature(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}

func TestFetchTempService_Success(t *testing.T) {
	mockRepo := new(MockTempRepository)
	service := usecase.NewFetchTempService(mockRepo)

	city := "São Paulo"
	expectedTemp := 25.0

	// Configuração do mock
	mockRepo.On("FetchTemperature", mock.Anything, city).Return(expectedTemp, nil)

	// Execução do teste
	temp, err := service.Fetch(context.Background(), city)

	// Validação
	require.NoError(t, err)
	require.Equal(t, expectedTemp, temp)

	mockRepo.AssertExpectations(t)
}

func TestFetchTempService_Error(t *testing.T) {
	mockRepo := new(MockTempRepository)
	service := usecase.NewFetchTempService(mockRepo)

	city := "Desconhecida"

	// Configuração do mock para erro
	mockRepo.On("FetchTemperature", mock.Anything, city).Return(0.0, errors.New("API error"))

	// Execução do teste
	temp, err := service.Fetch(context.Background(), city)

	// Validação
	require.Error(t, err)
	require.Zero(t, temp)

	mockRepo.AssertExpectations(t)
}
