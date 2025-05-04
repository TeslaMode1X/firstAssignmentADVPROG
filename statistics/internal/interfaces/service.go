package interfaces

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
)

type StatisticsService interface {
	GetInventoryStatisticsService() (*model.InventoryStatistics, error)
	GetOrdersStatisticsService() (*model.OrderStatistics, error)
	RecordProductActivity(p *model.Product) error
	RecordOrderActivity(order *model.Order) error
}
