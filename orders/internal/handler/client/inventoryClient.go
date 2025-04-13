package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ProductResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	StockLevel  int     `json:"stock"`
	CategoryID  string  `json:"category"`
}

type InventoryClient interface {
	GetProductByID(productID int64) (*ProductResponse, error)
	UpdateProductStock(productID int64, quantity int) error
}

type inventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewInventoryClient(baseURL string) InventoryClient {
	return &inventoryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *inventoryClient) GetProductByID(productID int64) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/product/%d", c.baseURL, productID)

	log.Printf("Sending GET request to: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Printf("Error during GET request: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Raw response from inventory service: %s", string(bodyBytes))

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("Product with ID %d not found (404)", productID)
		return nil, fmt.Errorf("product with ID %d not found", productID)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product ProductResponse
	if err := json.Unmarshal(bodyBytes, &product); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if product.ID == 0 {
		log.Printf("Invalid product: ID is 0")
		return nil, fmt.Errorf("invalid product: ID is 0")
	}

	log.Printf("Parsed product: ID=%d, Name=%s, StockLevel=%d, CategoryID=%s",
		product.ID, product.Name, product.StockLevel, product.CategoryID)

	return &product, nil
}

func (c *inventoryClient) UpdateProductStock(productID int64, quantity int) error {
	product, err := c.GetProductByID(productID)
	if err != nil {
		return err
	}

	if product.StockLevel < quantity {
		return fmt.Errorf("insufficient stock for product %d: have %d, need %d",
			productID, product.StockLevel, quantity)
	}

	url := fmt.Sprintf("%s/product", c.baseURL)

	categoryID := product.CategoryID
	if categoryID == "" {
		log.Printf("Warning: CategoryID is empty for product %d", productID)
		categoryID = "0"
	}

	updateData := map[string]interface{}{
		"id":          productID,
		"stock":       product.StockLevel - quantity,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"category":    categoryID,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Update product failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to update product stock, status code: %d", resp.StatusCode)
	}

	return nil
}
