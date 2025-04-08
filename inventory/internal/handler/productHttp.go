package handler

import (
	handler "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/handler/response"
	"net/http"
	"strconv"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ProductHttpHandler struct {
	productUsecase   usecase.ProductUsecase
	promotionUsecase usecase.PromotionUsecase
}

func NewProductHttpHandler(pu usecase.ProductUsecase, promotionUsecase usecase.PromotionUsecase) *ProductHttpHandler {
	return &ProductHttpHandler{productUsecase: pu, promotionUsecase: promotionUsecase}
}

func (p *ProductHttpHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.productUsecase.ProductDataProcessing(&product); err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, "success")
}

func (p *ProductHttpHandler) GetProducts(c *gin.Context) {
	products, err := p.productUsecase.ProductDataGetAll()
	if err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, products)
}

func (p *ProductHttpHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handler.Response(c, http.StatusBadRequest, "Invalid product ID: "+err.Error())
		return
	}

	product, err := p.productUsecase.ProductDataGetByID(id)
	if err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}
	if product == nil {
		handler.Response(c, http.StatusNotFound, "Product not found")
		return
	}

	handler.Response(c, http.StatusOK, product)
}

func (p *ProductHttpHandler) UpdateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		handler.Response(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := p.productUsecase.ProductDataUpdate(&product); err != nil {
		handler.Response(c, http.StatusInternalServerError, "Failed to update product: "+err.Error())
		return
	}

	handler.Response(c, http.StatusOK, product)
}

func (p *ProductHttpHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handler.Response(c, http.StatusBadRequest, "Invalid product ID: "+err.Error())
		return
	}

	if err := p.productUsecase.ProductDataDelete(id); err != nil {
		handler.Response(c, http.StatusInternalServerError, "Failed to delete product: "+err.Error())
		return
	}

	handler.Response(c, http.StatusOK, "success")
}

func (p *ProductHttpHandler) CreatePromotion(c *gin.Context) {
	var product model.Promotion
	if err := c.ShouldBindJSON(&product); err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.promotionUsecase.ProductDataCreatePromotion(&product); err != nil {
		handler.Response(c, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	handler.Response(c, http.StatusOK, product)
}

func (p *ProductHttpHandler) GetProductWithPromotion(c *gin.Context) {
	prods, err := p.promotionUsecase.ProductDataGetPromotions()
	if err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, prods)
}

func (p *ProductHttpHandler) DeletePromotion(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		handler.Response(c, http.StatusBadRequest, "Invalid promotion ID")
		return
	}

	if err := p.promotionUsecase.ProductDataDeletePromotion(idStr); err != nil {
		handler.Response(c, http.StatusInternalServerError, "Failed to delete promotion: "+err.Error())
		return
	}

	handler.Response(c, http.StatusOK, "success")
}
