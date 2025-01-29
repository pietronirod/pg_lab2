package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ViaCEPAPIURL  string `mapstructure:"VIACEP_API_URL"`
	WeatherAPIURL string `mapstructure:"WEATHERAPI_URL"`
	WeatherAPIKey string `mapstructure:"WEATHERAPI_KEY"`
	OTLPEndpoint  string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTLPProtocol  string `mapstructure:"OTEL_EXPORTER_OTLP_PROTOCOL"`
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	if AppConfig != nil {
		return AppConfig, nil
	}

	viper.SetConfigFile(".env") // Tenta carregar o arquivo .env
	viper.AutomaticEnv()        // Permite uso de variáveis de ambiente

	// Definir valores padrão para evitar falhas
	viper.SetDefault("VIACEP_API_URL", "https://viacep.com.br/ws/")
	viper.SetDefault("WEATHERAPI_URL", "http://api.weatherapi.com/v1/current.json")
	viper.SetDefault("WEATHERAPI_KEY", "")
	viper.SetDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317")
	viper.SetDefault("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")

	// Tentar ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Validação das configurações obrigatórias
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	AppConfig = &config
	return AppConfig, nil
}

func validateConfig(config *Config) error {
	if config.ViaCEPAPIURL == "" {
		return fmt.Errorf("VIACEP_API_URL is required")
	}
	if config.WeatherAPIURL == "" {
		return fmt.Errorf("WEATHERAPI_URL is required")
	}
	if config.WeatherAPIKey == "" {
		return fmt.Errorf("WEATHERAPI_KEY is required")
	}
	if config.OTLPEndpoint == "" {
		return fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is required")
	}
	if config.OTLPProtocol == "" {
		return fmt.Errorf("OTEL_EXPORTER_OTLP_PROTOCOL is required")
	}
	return nil
}
