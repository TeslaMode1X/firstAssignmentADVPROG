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
	productUsecase usecase.ProductUsecase
}

func NewProductHttpHandler(pu usecase.ProductUsecase) *ProductHttpHandler {
	return &ProductHttpHandler{productUsecase: pu}
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
