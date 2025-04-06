package dto

import (
	"time"
)

type ProductDTO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Price       float32   `gorm:"not null"`
	StockLevel  int       `gorm:"not null"`
	CategoryID  int64     `gorm:"index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   time.Time `gorm:"autoDeleteTime"`
}

type ProductResponse struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float32    `json:"price"`
	StockLevel  int        `json:"stock_level"`
	CategoryID  int64      `json:"category_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
