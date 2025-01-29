package delivery

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCEPHandler_Success(t *testing.T) {
	serviceBURL := "http://service-b:8090"
	mockClient := new(MockHTTPClient)
	handler := NewCEPHandler(serviceBURL, mockClient)

	cep := "01001000"
	requestBody := `{"cep":"` + cep + `"}`

	responseBody := `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, responseBody, w.Body.String())

	mockClient.AssertExpectations(t)
}

func TestCEPHandler_InvalidCEP(t *testing.T) {
	serviceBURL := "http://service-b:8090"
	mockClient := new(MockHTTPClient)
	handler := NewCEPHandler(serviceBURL, mockClient)

	cep := "123"
	requestBody := `{"cep":"` + cep + `"}`

	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	expectedResponse := `{"error":"invalid zipcode"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCEPHandler_CEPNotFound(t *testing.T) {
	serviceBURL := "http://service-b:8090"
	mockClient := new(MockHTTPClient)
	handler := NewCEPHandler(serviceBURL, mockClient)

	cep := "99999999"
	requestBody := `{"cep":"` + cep + `"}`

	responseBody := `{"error":"can not find zipcode"}`
	mockResponse := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, responseBody, w.Body.String())

	mockClient.AssertExpectations(t)
}

func TestCEPHandler_FetchCityError(t *testing.T) {
	serviceBURL := "http://service-b:8090"
	mockClient := new(MockHTTPClient)
	handler := NewCEPHandler(serviceBURL, mockClient)

	cep := "01001000"
	requestBody := `{"cep":"` + cep + `"}`

	responseBody := `{"error":"error fetching city"}`
	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, responseBody, w.Body.String())

	mockClient.AssertExpectations(t)
}

func TestCEPHandler_FetchTempError(t *testing.T) {
	serviceBURL := "http://service-b:8090"
	mockClient := new(MockHTTPClient)
	handler := NewCEPHandler(serviceBURL, mockClient)

	cep := "01001000"
	requestBody := `{"cep":"` + cep + `"}`

	responseBody := `{"error":"error fetching temperature"}`
	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, responseBody, w.Body.String())

	mockClient.AssertExpectations(t)
}
