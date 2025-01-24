package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pietronirod/lab2/service-a/otel"
)

var ErrNotFound = errors.New("location not found")

func GetLocationByCEP(ctx context.Context, cep string) (string, error) {
	_, span := otel.StartSpan(ctx, "ViaCEP Request")
	defer span.End()

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", ErrNotFound
	}
	defer resp.Body.Close()

	var data struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Localidade == "" {
		return "", ErrNotFound
	}
	return data.Localidade, nil
}
