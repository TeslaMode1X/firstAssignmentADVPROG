package handler

import (
	handler "github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler/response"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderHttpHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHttpHandler(orderUsecase usecase.OrderUsecase) *OrderHttpHandler {
	return &OrderHttpHandler{orderUsecase: orderUsecase}
}

func (o *OrderHttpHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := o.orderUsecase.CreateDataOrder(&order); err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, order)
}

func (o *OrderHttpHandler) GetOrders(c *gin.Context) {
	orders, err := o.orderUsecase.GetDataOrders()
	if err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, orders)
}

func (o *OrderHttpHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := o.orderUsecase.GetDataOrderByID(id)
	if err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, order)
}

func (o *OrderHttpHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	var orderMessage model.OrderMessage
	if err := c.ShouldBindJSON(&orderMessage); err != nil {
		handler.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := o.orderUsecase.UpdateDataOrderStatusByID(id, orderMessage.Message); err != nil {
		handler.Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Response(c, http.StatusOK, orderMessage)
}
