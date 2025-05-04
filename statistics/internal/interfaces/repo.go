package interfaces

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
)

type StatisticsRepo interface {
	GetInventoryStatisticsRepo() (*model.InventoryStatistics, error)
	GetOrdersStatisticsRepo() (*model.OrderStatistics, error)
	RecordProductActivityRepo(p *model.Product) error
	RecordOrderActivityRepo(order *model.Order) error
}
