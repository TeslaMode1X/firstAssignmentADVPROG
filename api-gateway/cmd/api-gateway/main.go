package main

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/api-gateway/internal/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	r := gin.Default()

	inventoryConn, err := grpc.NewClient("inventory:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to inventory: %v", err)
	}
	defer inventoryConn.Close()

	orderConn, err := grpc.NewClient("orders:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order: %v", err)
	}
	defer orderConn.Close()

	statisticsConn, err := grpc.NewClient("statistics:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to statistics: %v", err)
	}
	defer statisticsConn.Close()

	gatewayHandler := handler.NewGatewayHandler(inventoryConn, orderConn, statisticsConn)

	r.POST("/orders", gatewayHandler.CreateOrder)

	r.POST("/promotion", gatewayHandler.CreatePromotion)

	r.GET("/get/promotion", gatewayHandler.GetProductWithPromotion)

	r.GET("/orders", gatewayHandler.GetOrders)

	r.GET("/orders/:id", gatewayHandler.GetOrderByID)

	r.GET("/product", gatewayHandler.GetProducts)

	r.GET("/product/:id", gatewayHandler.GetProductByID)

	r.POST("/product", gatewayHandler.CreateProduct)

	r.PUT("/product", gatewayHandler.UpdateProduct)

	r.DELETE("/product/:id", gatewayHandler.DeleteProduct)

	r.DELETE("/promotion/:id", gatewayHandler.DeletePromotion)

	r.PATCH("/orders/:id", gatewayHandler.UpdateOrderStatus)

	r.GET("/statistics/inventory", gatewayHandler.GetInventoryStatistics)

	r.GET("/statistics/order", gatewayHandler.GetOrdersStatistics)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
