package main

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/api-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ordersURL := "http://orders_microservice:8082"
	inventoryURL := "http://inventory_microservice:8081"
	gatewayHandler := handler.NewGatewayHandler(ordersURL, inventoryURL)

	r.POST("/orders", gatewayHandler.CreateOrder)
	r.POST("/promotion", gatewayHandler.CreatePromotion)
	r.GET("/get/promotion", gatewayHandler.GetPromotions)

	r.GET("/orders", gatewayHandler.GetOrders)
	r.GET("/orders/:id", gatewayHandler.GetOrderByID)

	r.GET("/product", gatewayHandler.GetProducts)
	r.GET("/product/:id", gatewayHandler.GetProductByID)

	r.DELETE("/product/:id", gatewayHandler.DeleteProduct)
	r.DELETE("/promotion/:id", gatewayHandler.DeletePromotion)
	r.PATCH("/orders/:id", gatewayHandler.UpdateOrderStatus)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
