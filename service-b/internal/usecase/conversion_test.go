package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{37, 98.6},
	}

	for _, test := range tests {
		result := CelsiusToFahrenheit(test.celsius)
		assert.Equal(t, test.expected, result)
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 273.15},
		{100, 373.15},
		{-273.15, 0},
		{37, 310.15},
	}

	for _, test := range tests {
		result := CelsiusToKelvin(test.celsius)
		assert.Equal(t, test.expected, result)
	}
}
