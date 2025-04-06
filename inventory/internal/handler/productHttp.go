package handler

import (
	handler "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/handler/response"
	"net/http"

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
