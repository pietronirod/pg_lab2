package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock do serviço B para retornar uma resposta simulada
func mockServiceBResponse(statusCode int, responseBody string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})
	return httptest.NewServer(handler)
}

func TestCEPHandler_Success(t *testing.T) {
	mockServer := mockServiceBResponse(http.StatusOK, `{"city": "São Paulo", "temp_C": 25.5, "temp_F": 77.9, "temp_K": 298.65}`)
	defer mockServer.Close()

	handler := NewCEPHandler(mockServer.URL)

	reqBody := []byte(`{"cep": "01001000"}`)
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"city": "São Paulo", "temp_C": 25.5, "temp_F": 77.9, "temp_K": 298.65}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCEPHandler_InvalidCEP(t *testing.T) {
	handler := NewCEPHandler("http://mock-service-b")

	reqBody := []byte(`{"cep": "123"}`) // CEP inválido
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedResponse := `{"message": "Invalid CEP format", "code": 400, "trace_id": "00000000000000000000000000000000"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCEPHandler_InvalidRequestFormat(t *testing.T) {
	handler := NewCEPHandler("http://mock-service-b")

	reqBody := []byte(`{"invalid": "format"}`) // Formato de requisição inválido
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedResponse := `{"message": "Invalid request format", "code": 400, "trace_id": "00000000000000000000000000000000"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCEPHandler_ServiceBUnavailable(t *testing.T) {
	// Criamos um servidor que retorna "Service B unavailable"
	mockServer := mockServiceBResponse(http.StatusInternalServerError, `{"message": "Service B unavailable"}`)
	defer mockServer.Close()

	handler := NewCEPHandler(mockServer.URL)

	reqBody := []byte(`{"cep": "01001000"}`)
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedResponse := `{"message": "Failed to contact service-b", "code": 500, "trace_id": "00000000000000000000000000000000"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
