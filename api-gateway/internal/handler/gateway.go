package handler

import (
	"bytes"
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/inventory"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"strconv"
)

type GatewayHandler struct {
	inventoryClient inventory.InventoryServiceClient
	orderClient     orders.OrderServiceClient
	httpClient      *http.Client
}

func NewGatewayHandler(inventoryConn *grpc.ClientConn, orderConn *grpc.ClientConn) *GatewayHandler {
	return &GatewayHandler{
		inventoryClient: inventory.NewInventoryServiceClient(inventoryConn),
		orderClient:     orders.NewOrderServiceClient(orderConn),
		httpClient:      &http.Client{},
	}
}

func (h *GatewayHandler) CreateProduct(c *gin.Context) {
	var req inventory.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.inventoryClient.CreateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetProduct())
}

func (h *GatewayHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.inventoryClient.GetProductByID(context.Background(), &inventory.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetProduct())
}

func (h *GatewayHandler) UpdateProduct(c *gin.Context) {
	var req inventory.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//idStr := c.Param("id")
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	log.Printf("Received update request: %+v", req)

	resp, err := h.inventoryClient.UpdateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetProduct())
}

func (h *GatewayHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := h.inventoryClient.DeleteProduct(context.Background(), &inventory.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *GatewayHandler) GetProducts(c *gin.Context) {
	var req inventory.Empty
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.inventoryClient.GetProducts(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetProducts())
}

//func (h *GatewayHandler) CreatePromotion(c *gin.Context) {
//	h.proxyRequest(c, h.ordersURL+"/product/promotion", "POST")
//}
//
//func (h *GatewayHandler) GetPromotions(c *gin.Context) {
//	h.proxyRequest(c, h.ordersURL+"/product/promotion", "GET")
//}

func (h *GatewayHandler) CreateOrder(c *gin.Context) {
	var order orders.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.orderClient.CreateOrder(c.Request.Context(), &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *GatewayHandler) GetOrders(c *gin.Context) {
	var req orders.Empty
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.orderClient.GetOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetOrder())
}

func (h *GatewayHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	resp, err := h.orderClient.GetOrderById(c.Request.Context(), &orders.GetOrderRequest{
		Id: int32(orderId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GatewayHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var requestBody struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.orderClient.UpdateOrderStatusById(c.Request.Context(), &orders.UpdateOrderStatusByIdRequest{
		Id:      int32(orderId),
		Message: requestBody.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

//func (h *GatewayHandler) DeletePromotion(c *gin.Context) {
//	id := c.Param("id")
//	h.proxyRequest(c, h.ordersURL+"/product/promotion/"+id, "DELETE")
//}

func (h *GatewayHandler) proxyRequest(c *gin.Context, url, method string) {
	body, _ := io.ReadAll(c.Request.Body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header = c.Request.Header

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
