package dto

import "time"

type PromotionDTO struct {
	ID                 string  `gorm:"primary_key"`
	Name               string  `gorm:"type:varchar(255);not null"`
	Description        string  `gorm:"type:varchar(255);not null"`
	DiscountPercentage float64 `gorm:"type:float;not null"`
	ApplicableProducts string  `gorm:"type:varchar(255);not null"`
	StartDate          time.Time
	EndDate            time.Time
	IsActive           bool `gorm:"type:boolean;not null"`
}

func (PromotionDTO) TableName() string {
	return "promotion"
}

type PromotionResponse struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ApplicableProducts []string  `json:"applicable_products"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	IsActive           bool      `json:"is_active"`
}
