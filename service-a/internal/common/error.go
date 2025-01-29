package common

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse define o formato padronizado para respostas de erro
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	TraceID string `json:"trace_id,omitempty"` // Campo opcional para correlação de tracing
}

// NewErrorResponse cria uma nova resposta de erro padronizada
func NewErrorResponse(w http.ResponseWriter, statusCode int, message string, traceID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
		Code:    statusCode,
		TraceID: traceID,
	}

	// Log do erro
	log.Printf("Error: %s, Code: %d, TraceID: %s", message, statusCode, traceID)

	json.NewEncoder(w).Encode(response)
}
