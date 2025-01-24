package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pietronirod/lab2/service-a/handler"
	"github.com/pietronirod/lab2/service-a/otel"
)

func main() {
	tracerShutdown := otel.InitTracer("service-a", "https://localhost:9411/api/v2/spans")
	defer tracerShutdown()

	otel.InitPropagator()

	router := gin.Default()
	router.POST("/input", handler.HandleInput)

	log.Println("Service A running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start Service A: %v", err)
	}
}
