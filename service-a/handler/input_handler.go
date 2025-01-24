package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pietronirod/lab2/service-a/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Input struct {
	CEP string `json:"cep" binding:"required,len=8"`
}

func HandleInput(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	ctx, span := otel.StartSpan(c.Request.Context(), "Send to Service B")
	defer span.End()

	body, _ := json.Marshal(input)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://service-b:8082/weather", bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error creating request"})
		return
	}

	otel.GetPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "service b unavailable"})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"data": resp.Body})
}
