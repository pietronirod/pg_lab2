package delivery

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// HTTPClient é uma interface para o cliente HTTP
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CEPHandler gerencia as requisições para enviar o CEP ao serviço B
type CEPHandler struct {
	serviceBURL string
	httpClient  HTTPClient
}

// NewCEPHandler cria um novo handler
func NewCEPHandler(serviceBURL string, httpClient HTTPClient) *CEPHandler {
	return &CEPHandler{serviceBURL: serviceBURL, httpClient: httpClient}
}

// Handle processa a requisição para enviar o CEP ao serviço B
func (h *CEPHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("CEPHandler: Request received")

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "process-cep-handler")
	defer span.End()
	log.Printf("CEPHandler: Received TraceID=%s", span.SpanContext().TraceID().String())

	var requestBody struct {
		CEP string `json:"cep"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("CEPHandler: Invalid request body: %v", err)
		span.SetStatus(codes.Error, "Invalid request body")
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validação do CEP
	if len(requestBody.CEP) != 8 {
		log.Printf("CEPHandler: Invalid CEP: %s", requestBody.CEP)
		span.SetStatus(codes.Error, "Invalid CEP length")
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}
	span.SetAttributes(attribute.String("cep", requestBody.CEP))

	// Enviar CEP ao serviço B
	serviceBURL := h.serviceBURL + "/cep/" + requestBody.CEP
	req, err := http.NewRequestWithContext(ctx, "GET", serviceBURL, nil)
	if err != nil {
		log.Printf("CEPHandler: Error creating request to service B: %v", err)
		span.SetStatus(codes.Error, "Error creating request to service B")
		h.writeErrorResponse(w, http.StatusInternalServerError, "error creating request to service B")
		return
	}

	// Propagar o contexto de rastreamento na requisição HTTP
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("CEPHandler: Error contacting service B: %v", err)
		span.SetStatus(codes.Error, "Error contacting service B")
		h.writeErrorResponse(w, http.StatusInternalServerError, "error contacting service B")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("CEPHandler: Error reading response from service B: %v", err)
		span.SetStatus(codes.Error, "Error reading response from service B")
		h.writeErrorResponse(w, http.StatusInternalServerError, "error reading response from service B")
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("CEPHandler: Service B returned error: %s", body)
		span.SetStatus(codes.Error, "Service B returned error")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Responder com a resposta do serviço B
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *CEPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + message + `"}`))
}
