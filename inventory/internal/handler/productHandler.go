package handler

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type PromotionHandler interface {
	CreatePromotion(c *gin.Context)
	GetProductWithPromotion(c *gin.Context)
	DeletePromotion(c *gin.Context)
}
