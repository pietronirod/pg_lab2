package common

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse define o formato padronizado para respostas de erro
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewErrorResponse cria uma nova resposta de erro padronizada
func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(response)
}
