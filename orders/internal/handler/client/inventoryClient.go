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
	ID          int64      `json:"ID"`
	Name        string     `json:"Name"`
	Description string     `json:"Description"`
	Price       float32    `json:"Price"`
	StockLevel  int        `json:"StockLevel"`
	CategoryID  int64      `json:"CategoryID"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
	DeletedAt   *time.Time `json:"DeletedAt"`
}

type ResponseWrapper struct {
	Message ProductResponse `json:"message"`
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

	var responseWrapper ResponseWrapper
	if err := json.Unmarshal(bodyBytes, &responseWrapper); err != nil {
		log.Printf("Failed to decode response with wrapper: %v", err)

		var product ProductResponse
		if err := json.Unmarshal(bodyBytes, &product); err != nil {
			log.Printf("Failed to decode response directly: %v", err)
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		log.Printf("Parsed product directly: ID=%d, Name=%s, StockLevel=%d",
			product.ID, product.Name, product.StockLevel)
		return &product, nil
	}

	product := responseWrapper.Message
	log.Printf("Parsed product from wrapper: ID=%d, Name=%s, StockLevel=%d",
		product.ID, product.Name, product.StockLevel)

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

	updateData := map[string]interface{}{
		"ID":          productID,
		"StockLevel":  product.StockLevel - quantity,
		"Name":        product.Name,
		"Description": product.Description,
		"Price":       product.Price,
		"CategoryID":  product.CategoryID,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
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
