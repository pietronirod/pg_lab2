package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pietronirod/lab2/service-a/otel"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetEnvPrefix("weather")
	viper.BindEnv("api_key")
	viper.SetDefault("api_key", "")
}

func GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	_, span := otel.StartSpan(ctx, "WeatherAPI Request")
	defer span.End()

	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return 0, fmt.Errorf("WEATHER_API_KEY is not set")
	}

	url := "http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, err
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}
