package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceBURL   string `mapstructure:"SERVICE_B_URL"`
	ViaCEPAPIURL  string `mapstructure:"VIACEP_API_URL"`
	WeatherAPIURL string `mapstructure:"WEATHERAPI_URL"`
	WeatherAPIKey string `mapstructure:"WEATHERAPI_KEY"`
	OTLPEndpoint  string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

var AppConfig *Config

func LoadConfig() *Config {
	if AppConfig != nil {
		return AppConfig
	}

	viper.SetConfigFile(".env") // Carrega o arquivo .env
	viper.AutomaticEnv()        // Permite usar variáveis de ambiente diretamente

	// Valores padrão (fallbacks)
	viper.SetDefault("SERVICE_B_URL", "http://service-b:8090")
	viper.SetDefault("VIACEP_API_URL", "https://viacep.com.br/ws/")
	viper.SetDefault("WEATHERAPI_URL", "http://api.weatherapi.com/v1/current.json")
	viper.SetDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4317")

	// Lê as configurações
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	AppConfig = &config
	return AppConfig
}
