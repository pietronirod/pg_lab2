package delivery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := tracer.Start(ctx, "process-cep-handler")
	defer span.End()

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CEPHandler: Error decoding request body: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid JSON body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("CEPHandler: CEP received %s", req.CEP)
	span.SetAttributes(attribute.String("cep", req.CEP))

	// Validar o CEP (precisa ter exatamente 8 dígitos)
	if len(req.CEP) != 8 {
		log.Printf("CEPHandler: Invalid CEP: %v", req.CEP)
		span.SetStatus(codes.Error, "Invalid CEP length")
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Montar URL para chamar o Service B
	serviceBURL := h.serviceBURL + "/cep/" + req.CEP
	log.Printf("CEPHandler: Calling Service B at URL: %s", serviceBURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", serviceBURL, nil)
	if err != nil {
		log.Printf("CEPHandler: Error creating request to Service B: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error creating request to Service B")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Instrumentar o cliente HTTP com OpenTelemetry
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("CEPHandler: Error calling Service B: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error calling Service B")
		http.Error(w, "Service B unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("CEPHandler: Response from Service B with status code: %d", resp.StatusCode)
	span.SetAttributes(attribute.Int("service-b.status_code", resp.StatusCode))

	// Encaminhar a resposta do Service B para o cliente
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("CEPHandler: Error reading response from Service B: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error reading response from Service B")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	span.SetStatus(codes.Ok, "Successfully processed CEP")
	log.Println("CEPHandler: Request processed successfully")
}
