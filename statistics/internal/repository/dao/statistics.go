package dao

import (
	"gorm.io/gorm"
)

type OrderStatistics struct {
	gorm.Model
	TotalOrders  int     `gorm:"column:total_orders"`
	TotalRevenue float32 `gorm:"column:total_revenue"`
	AveragePrice float32 `gorm:"column:average_price"`
}

type InventoryStatistics struct {
	gorm.Model
	TotalProducts       int     `gorm:"column:total_products"`
	TotalStock          int     `gorm:"column:total_stock"`
	TotalInventoryValue float32 `gorm:"column:total_inventory_value"`
}
