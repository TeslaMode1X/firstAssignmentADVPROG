package model

import (
	"fmt"
	"time"
)

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

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be positive")
	}
	if p.StockLevel < 0 {
		return fmt.Errorf("stock level cannot be negative")
	}
	return nil
}
