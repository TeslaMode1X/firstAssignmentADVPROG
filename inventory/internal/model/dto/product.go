package dto

import (
	"gorm.io/gorm"
	"time"
)

type ProductDTO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       float32        `gorm:"not null"`
	StockLevel  int            `gorm:"not null"`
	CategoryID  int64          `gorm:"index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"autoDeleteTime"`
}

func (CategoryDTO) TableName() string {
	return "products"
}

type ProductResponse struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ApplicableProducts []string  `json:"applicable_products"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	IsActive           bool      `json:"is_active"`
}
