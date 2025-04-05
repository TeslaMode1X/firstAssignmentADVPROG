package entity

import "time"

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
