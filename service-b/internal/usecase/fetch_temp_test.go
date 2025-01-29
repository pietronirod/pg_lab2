package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTemperatureRepository struct {
	mock.Mock
}

func (m *MockTemperatureRepository) FetchTemperature(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}

func TestFetchTempService_Success(t *testing.T) {
	mockRepo := new(MockTemperatureRepository)
	service := NewFetchTempService(mockRepo)

	city := "São Paulo"
	expectedTemp := 25.5

	// Configuração do mock
	mockRepo.On("FetchTemperature", mock.Anything, city).Return(expectedTemp, nil)

	// Execução do teste
	temp, err := service.Fetch(context.Background(), city)

	// Validação
	require.NoError(t, err)
	require.Equal(t, expectedTemp, temp)

	mockRepo.AssertExpectations(t)
}

func TestFetchTempService_APIError(t *testing.T) {
	mockRepo := new(MockTemperatureRepository)
	service := NewFetchTempService(mockRepo)

	city := "São Paulo"
	expectedError := fmt.Errorf("API error")

	// Configuração do mock para erro de comunicação com a API
	mockRepo.On("FetchTemperature", mock.Anything, city).Return(0.0, expectedError)

	// Execução do teste
	temp, err := service.Fetch(context.Background(), city)

	// Validação
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.Equal(t, 0.0, temp)

	mockRepo.AssertExpectations(t)
}
