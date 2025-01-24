package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pietronirod/lab2/service-a/otel"
	"github.com/pietronirod/lab2/service-b/service"
	"go.opentelemetry.io/otel/propagation"
)

func HandleWeather(c *gin.Context) {
	ctx := otel.GetPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

	ctx, span := otel.StartSpan(ctx, "Process Weather Request")
	defer span.End()

	var input struct {
		CEP string `json:"cep"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || len(input.CEP) != 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	city, err := service.GetLocationByCEP(ctx, input.CEP)
	if errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "cannot find zipcode"})
		return
	}

	tempC, err := service.GetTemperatureByCity(ctx, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "temperature service error"})
		return
	}

	response := map[string]interface{}{
		"city":   city,
		"temp_C": tempC,
		"temp_F": tempC*1.8 + 32,
		"temp_K": tempC + 273.15,
	}
	c.JSON(http.StatusOK, response)
}
