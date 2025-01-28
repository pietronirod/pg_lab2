package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CEPService interface {
	GetCEPInfo(ctx context.Context, cep string) (interface{}, error)
}

type CEPServiceImpl struct{}

func NewCEPService() CEPService {
	return &CEPServiceImpl{}
}

func (s *CEPServiceImpl) GetCEPInfo(ctx context.Context, cep string) (interface{}, error) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "cep-service")
	defer span.End()

	serviceBURL := "http://service-b:8090/cep/" + cep
	req, _ := http.NewRequestWithContext(ctx, "GET", serviceBURL, nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("CEPService: Error calling Service B: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("CEPService: Service B returned status: %d", resp.StatusCode)
		return nil, errors.New("service B error")
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
