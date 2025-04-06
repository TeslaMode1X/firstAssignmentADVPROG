package model

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          int64
	Items       []OrderItem
	TotalAmount float32
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

type OrderMessage struct {
	Message string
}

type OrderItem struct {
	ProductID   int64   `json:"product_id"`
	Quantity    int     `json:"quantity"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
}

func (o *Order) Validate() error {
	if len(o.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}
	for _, item := range o.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("product ID must be positive")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item quantity must be positive")
		}
	}
	return nil
}
