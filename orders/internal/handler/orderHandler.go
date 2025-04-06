package handler

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
}
