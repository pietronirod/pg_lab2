package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type CEPService interface {
	GetCEPInfo(ctx context.Context, cep string) (interface{}, error)
}

type CEPServiceImpl struct {
	serviceBURL string
}

func NewCEPService(serviceBURL string) CEPService {
	return &CEPServiceImpl{serviceBURL: serviceBURL}
}

func (s *CEPServiceImpl) GetCEPInfo(ctx context.Context, cep string) (interface{}, error) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "cep-service")
	defer span.End()

	url := s.serviceBURL + "/cep/" + cep
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("CEPService: Error calling Service B: %v", err)
		span.SetStatus(codes.Error, "Error calling Service B")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("CEPService: Service B returned status: %d", resp.StatusCode)
		span.SetStatus(codes.Error, "Service B returned non-OK status")
		return nil, errors.New("service B error")
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		span.SetStatus(codes.Error, "Error decoding response from Service B")
		return nil, err
	}

	span.SetStatus(codes.Ok, "Successfully retrieved CEP info")
	return result, nil
}
