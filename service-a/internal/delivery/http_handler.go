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
	"go.opentelemetry.io/otel/propagation"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func CEPHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(req.CEP) != 8 {
		log.Printf("CEPHandler: Invalid CEP: %v", req.CEP)
		span.SetStatus(codes.Error, "Invalid CEP length")
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	serviceBURL := "http://service-b:8090/cep/" + req.CEP
	log.Printf("CEPHandler: Calling Service B at URL: %s", serviceBURL)
	log.Printf("CEPHandler: Propagating TraceID=%s to Service B", span.SpanContext().TraceID().String())

	httpReq, err := http.NewRequestWithContext(ctx, "GET", serviceBURL, nil)
	if err != nil {
		log.Printf("CEPHandler: Error creating request to Service B: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error creating request to Service B")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("CEPHandler: Error calling Service B: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error calling Service B")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("CEPHandler: Response from Service B with status code: %d", resp.StatusCode)
	span.SetAttributes(attribute.Int("service-b.status_code", resp.StatusCode))

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
