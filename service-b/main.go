package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pietronirod/lab2/service-a/otel"
	"github.com/pietronirod/lab2/service-b/handler"
	"github.com/pietronirod/lab2/service-b/service"
)

func main() {
	service.InitConfig()

	tracerShutdown := otel.InitTracer("service-b", "http://localhost:9411/api/v2/spans")
	defer tracerShutdown()

	router := gin.Default()
	router.POST("/weather", handler.HandleWeather)

	log.Println("Service B running on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Failed to start Service B: %v", err)
	}
}
