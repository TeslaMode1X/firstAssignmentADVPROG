package dto

import (
	"time"
)

type OrderDTO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	UserID      int64          `gorm:"index;not null"`
	Status      string         `gorm:"type:varchar(50);not null"`
	TotalAmount float32        `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	Items       []OrderItemDTO `gorm:"foreignKey:OrderID"`
}

func (OrderDTO) TableName() string {
	return "orders"
}

type OrderItemDTO struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	OrderID     int64   `gorm:"index;not null"`
	ProductID   int64   `gorm:"index;not null"`
	Quantity    int     `gorm:"not null"`
	Price       float32 `gorm:"not null"`
	ProductName string  `gorm:"type:varchar(255)"`
}

func (OrderItemDTO) TableName() string {
	return "order_items"
}

type OrderResponse struct {
	ID          int64               `json:"id"`
	UserID      int64               `json:"user_id"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float32             `json:"total_amount"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
}

type OrderItemResponse struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
	Subtotal    float32 `json:"subtotal"`
}
