package main

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/api-gateway/internal/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
)

func main() {
	r := gin.Default()

	inventoryConn, err := grpc.Dial("inventory:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory: %v", err)
	}
	defer inventoryConn.Close()

	ordersURL := "http://orders_microservice:8082"

	gatewayHandler := handler.NewGatewayHandler(inventoryConn, ordersURL)

	r.POST("/orders", gatewayHandler.CreateOrder)
	r.POST("/promotion", gatewayHandler.CreatePromotion)

	r.GET("/get/promotion", gatewayHandler.GetPromotions)

	r.GET("/orders", gatewayHandler.GetOrders)
	r.GET("/orders/:id", gatewayHandler.GetOrderByID)

	r.GET("/product", gatewayHandler.GetProducts)
	r.GET("/product/:id", gatewayHandler.GetProductByID)

	r.POST("/product", gatewayHandler.CreateProduct)

	r.PUT("/product", gatewayHandler.UpdateProduct)

	r.DELETE("/product/:id", gatewayHandler.DeleteProduct)
	r.DELETE("/promotion/:id", gatewayHandler.DeletePromotion)
	r.PATCH("/orders/:id", gatewayHandler.UpdateOrderStatus)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
