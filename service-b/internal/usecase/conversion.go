package usecase

import "math"

// CelsiusToFahrenheit converte Celsius para Fahrenheit e arredonda para 2 casas decimais
func CelsiusToFahrenheit(c float64) float64 {
	return math.Round((c*1.8+32)*100) / 100
}

// CelsiusToKelvin converte Celsius para Kelvin e arredonda para 2 casas decimais
func CelsiusToKelvin(c float64) float64 {
	return math.Round((c+273.15)*100) / 100
}
