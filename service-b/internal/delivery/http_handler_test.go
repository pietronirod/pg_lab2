package delivery

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"service-b/internal/repository"
	"service-b/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFetchCityService struct {
	mock.Mock
}

func (m *MockFetchCityService) Fetch(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}

type MockFetchTempService struct {
	mock.Mock
}

func (m *MockFetchTempService) Fetch(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}

func TestCEPHandler_Success(t *testing.T) {
	mockFetchCity := new(MockFetchCityService)
	mockFetchTemp := new(MockFetchTempService)
	handler := NewCEPHandler(mockFetchCity, mockFetchTemp)

	cep := "01001000"
	expectedCity := "São Paulo"
	expectedTempC := 28.5
	expectedTempF := usecase.CelsiusToFahrenheit(expectedTempC)
	expectedTempK := usecase.CelsiusToKelvin(expectedTempC)

	// Configuração dos mocks
	mockFetchCity.On("Fetch", mock.Anything, cep).Return(expectedCity, nil)
	mockFetchTemp.On("Fetch", mock.Anything, expectedCity).Return(expectedTempC, nil)

	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep, nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"city":"São Paulo","temp_C":28.5,"temp_F":` + fmt.Sprintf("%.2f", expectedTempF) + `,"temp_K":` + fmt.Sprintf("%.2f", expectedTempK) + `}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockFetchCity.AssertExpectations(t)
	mockFetchTemp.AssertExpectations(t)
}

func TestCEPHandler_InvalidCEP(t *testing.T) {
	mockFetchCity := new(MockFetchCityService)
	mockFetchTemp := new(MockFetchTempService)
	handler := NewCEPHandler(mockFetchCity, mockFetchTemp)

	cep := "123"

	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep, nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	expectedResponse := `{"error":"invalid zipcode"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCEPHandler_CEPNotFound(t *testing.T) {
	mockFetchCity := new(MockFetchCityService)
	mockFetchTemp := new(MockFetchTempService)
	handler := NewCEPHandler(mockFetchCity, mockFetchTemp)

	cep := "99999999"

	// Configuração do mock para erro de CEP não encontrado
	mockFetchCity.On("Fetch", mock.Anything, cep).Return("", repository.ErrCEPNotFound)

	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep, nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	expectedResponse := `{"error":"can not find zipcode"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockFetchCity.AssertExpectations(t)
}

func TestCEPHandler_FetchCityError(t *testing.T) {
	mockFetchCity := new(MockFetchCityService)
	mockFetchTemp := new(MockFetchTempService)
	handler := NewCEPHandler(mockFetchCity, mockFetchTemp)

	cep := "01001000"
	expectedError := fmt.Errorf("city not found")

	// Configuração do mock para erro ao buscar cidade
	mockFetchCity.On("Fetch", mock.Anything, cep).Return("", expectedError)

	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep, nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedResponse := `{"error":"error fetching city"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockFetchCity.AssertExpectations(t)
}

func TestCEPHandler_FetchTempError(t *testing.T) {
	mockFetchCity := new(MockFetchCityService)
	mockFetchTemp := new(MockFetchTempService)
	handler := NewCEPHandler(mockFetchCity, mockFetchTemp)

	cep := "01001000"
	expectedCity := "São Paulo"
	expectedError := fmt.Errorf("temperature not found")

	// Configuração dos mocks
	mockFetchCity.On("Fetch", mock.Anything, cep).Return(expectedCity, nil)
	mockFetchTemp.On("Fetch", mock.Anything, expectedCity).Return(0.0, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/cep/"+cep, nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedResponse := `{"error":"error fetching temperature"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockFetchCity.AssertExpectations(t)
	mockFetchTemp.AssertExpectations(t)
}
