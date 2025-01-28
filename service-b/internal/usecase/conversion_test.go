package usecase_test

import (
	"service-b/internal/usecase"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	require.Equal(t, 32.0, usecase.CelsiusToFahrenheit(0))
	require.Equal(t, 212.0, usecase.CelsiusToFahrenheit(100))
	require.Equal(t, 98.6, usecase.CelsiusToFahrenheit(37))
}

func TestCelsiusToKelvin(t *testing.T) {
	require.Equal(t, 273.15, usecase.CelsiusToKelvin(0))
	require.Equal(t, 373.15, usecase.CelsiusToKelvin(100))
	require.Equal(t, 310.15, usecase.CelsiusToKelvin(37))
}
