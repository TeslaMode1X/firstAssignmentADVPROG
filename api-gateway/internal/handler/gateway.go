package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type GatewayHandler struct {
	ordersURL    string
	inventoryURL string
	client       *http.Client
}

func NewGatewayHandler(ordersURL, inventoryURL string) *GatewayHandler {
	return &GatewayHandler{
		ordersURL:    ordersURL,
		inventoryURL: inventoryURL,
		client:       &http.Client{},
	}
}

func (h *GatewayHandler) CreateOrder(c *gin.Context) {
	h.proxyRequest(c, h.ordersURL+"/order/create", "POST")
}

func (h *GatewayHandler) GetOrders(c *gin.Context) {
	h.proxyRequest(c, h.ordersURL+"/order", "GET")
}

func (h *GatewayHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, h.ordersURL+"/order/"+id, "GET")
}

func (h *GatewayHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, h.ordersURL+"/order/"+id, "PATCH")
}

func (h *GatewayHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, h.inventoryURL+"/product/"+id, "GET")
}

func (h *GatewayHandler) GetProducts(c *gin.Context) {
	h.proxyRequest(c, h.inventoryURL+"/product", "GET")
}

func (h *GatewayHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, h.inventoryURL+"/product/"+id, "DELETE")
}

func (h *GatewayHandler) proxyRequest(c *gin.Context, url, method string) {
	body, _ := io.ReadAll(c.Request.Body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header = c.Request.Header

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
