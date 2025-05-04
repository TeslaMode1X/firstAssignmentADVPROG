package model

import "time"

type OrderStatistics struct {
	ID           int
	TotalOrders  int
	TotalRevenue float32
	AveragePrice float32
}

type InventoryStatistics struct {
	ID                  int
	TotalProducts       int
	TotalStock          int
	TotalInventoryValue int
}

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float32
	StockLevel  int
	CategoryID  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Order struct {
	ID          int64
	Items       []OrderItem
	TotalAmount float32
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

type OrderStatus string

type OrderItem struct {
	ProductID   int64   `json:"product_id"`
	Quantity    int     `json:"quantity"`
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
}
