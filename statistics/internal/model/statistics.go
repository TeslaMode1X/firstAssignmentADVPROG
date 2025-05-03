package model

import "time"

type OrderStatistics struct {
	TotalOrders  int
	TotalRevenue float32
	AveragePrice float32
}

type InventoryStatistics struct {
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
