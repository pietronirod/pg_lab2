package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCityFromCEP_Success(t *testing.T) {
	mockResponse := `{"localidade":"São Paulo"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalURL := "https://viacep.com.br/ws/"
	defer func() { apiBaseURL = originalURL }()
	apiBaseURL = server.URL + "/"

	city, err := FetchCityFromCEP("12345678")
	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
}

func TestFetchCityFromCEP_NotFound(t *testing.T) {
	mockResponse := `{"erro": true}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalURL := "https://viacep.com.br/ws/"
	defer func() { apiBaseURL = originalURL }()

	city, err := FetchCityFromCEP("00000000")
	assert.Error(t, err)
	assert.Equal(t, ErrCEPNotFound, err)
	assert.Equal(t, "", city)
}
