package delivery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"service-a/internal/common"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// CEPRequest define o formato esperado da requisição
type CEPRequest struct {
	CEP string `json:"cep"`
}

// CEPHandler gerencia as requisições do service-a
type CEPHandler struct {
	serviceBURL string
}

// NewCEPHandler cria um novo CEPHandler
func NewCEPHandler(serviceBURL string) *CEPHandler {
	return &CEPHandler{serviceBURL: serviceBURL}
}

// Handle processa a requisição para encaminhar ao service-b
func (h *CEPHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("CEPHandler: Request received")

	ctx := r.Context()
	tracer := otel.Tracer("service-a")
	_, span := tracer.Start(ctx, "HandleRequest")
	defer span.End()

	// Decodificar a requisição
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CEP == "" {
		span.SetStatus(codes.Error, "Invalid request format")
		common.NewErrorResponse(w, http.StatusBadRequest, "Invalid request format", span.SpanContext().TraceID().String())
		return
	}

	// Validação do CEP
	if len(req.CEP) != 8 {
		span.SetStatus(codes.Error, "Invalid CEP format")
		common.NewErrorResponse(w, http.StatusBadRequest, "Invalid CEP format", span.SpanContext().TraceID().String())
		return
	}

	// Encaminhar a requisição ao service-b
	url := h.serviceBURL + "/cep/" + req.CEP

	// Criar uma nova requisição HTTP com o contexto de tracing
	reqServiceB, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create request to service-b")
		common.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create request to service-b", span.SpanContext().TraceID().String())
		return
	}

	// Propagar o contexto de tracing na requisição
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(reqServiceB.Header))

	// Enviar a requisição ao service-b
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(reqServiceB)
	if err != nil || resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "Failed to contact service-b")
		common.NewErrorResponse(w, http.StatusInternalServerError, "Failed to contact service-b", span.SpanContext().TraceID().String())
		if err == nil {
			defer resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()

	// Ler a resposta do service-b
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to read response from service-b")
		common.NewErrorResponse(w, http.StatusInternalServerError, "Failed to read response from service-b", span.SpanContext().TraceID().String())
		return
	}

	// Encaminhar a resposta ao cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	span.SetStatus(codes.Ok, "Request handled successfully")
	span.SetAttributes(attribute.String("CEP", req.CEP))
}
