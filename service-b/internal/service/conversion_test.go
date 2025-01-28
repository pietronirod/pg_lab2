package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	assert.InDelta(t, 32.0, CelsiusToFahrenheit(0), 0.0001)
	assert.InDelta(t, 212.0, CelsiusToFahrenheit(100), 0.0001)
	assert.InDelta(t, 98.6, CelsiusToFahrenheit(37), 0.0001)
}

func TestCelsiusToKelvin(t *testing.T) {
	assert.InDelta(t, 273.15, CelsiusToKelvin(0), 0.0001)
	assert.InDelta(t, 373.15, CelsiusToKelvin(100), 0.0001)
	assert.InDelta(t, 310.15, CelsiusToKelvin(37), 0.0001)
}
