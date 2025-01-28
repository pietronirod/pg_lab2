package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCEPHandler_InvalidCEP(t *testing.T) {
	body := []byte(`{"cep":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(body))
	w := httptest.NewRecorder()

	CEPHandler(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "invalid zipcode\n", w.Body.String())
}

func TestCEPHandler_ValidCEP(t *testing.T) {
	body := []byte(`{"cep":"12345678"}`)
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(body))
	w := httptest.NewRecorder()

	callServiceB = func(_ string) ([]byte, error) {
		return []byte(`{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`), nil
	}

	CEPHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
