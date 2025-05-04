package repository

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
)

type StatisticsRepo struct {
	db interfaces.Database
}

func NewStatisticsRepo(db interfaces.Database) interfaces.StatisticsRepo {
	return &StatisticsRepo{
		db: db,
	}
}

func (r *StatisticsRepo) GetInventoryStatisticsRepo() (*model.InventoryStatistics, error) {
	const op = "repository.statistics.GetInventoryStatisticsRepo"

	var stats model.InventoryStatistics
	result := r.db.GetDB().Model(&model.InventoryStatistics{}).First(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stats, nil
}

func (r *StatisticsRepo) GetOrdersStatisticsRepo() (*model.OrderStatistics, error) {
	const op = "repository.statistics.GetOrdersStatisticsRepo"

	var stats model.OrderStatistics
	result := r.db.GetDB().Model(&model.OrderStatistics{}).First(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stats, nil
}

func (r *StatisticsRepo) RecordProductActivityRepo(p *model.Product) error {
	const op = "repository.statistics.RecordProductActivityRepo"

	var stats model.InventoryStatistics
	result := r.db.GetDB().FirstOrCreate(&stats, model.InventoryStatistics{ID: 1})
	if result.Error != nil {
		return result.Error
	}

	stats.TotalProducts += 1
	stats.TotalStock += p.StockLevel
	stats.TotalInventoryValue += int(p.Price) * p.StockLevel

	if err := r.db.GetDB().Save(&stats).Error; err != nil {
		return err
	}

	return nil
}

func (r *StatisticsRepo) RecordOrderActivityRepo(order *model.Order) error {
	const op = "repository.statistics.RecordOrderActivityRepo"

	var stats model.OrderStatistics
	result := r.db.GetDB().FirstOrCreate(&stats, model.OrderStatistics{ID: 1})
	if result.Error != nil {
		return fmt.Errorf("%s: failed to get or create stats: %w", op, result.Error)
	}

	stats.TotalOrders += 1

	for _, item := range order.Items {
		stats.TotalRevenue += item.Price
	}

	if stats.TotalOrders > 0 {
		stats.AveragePrice = stats.TotalRevenue / float32(stats.TotalOrders)
	} else {
		stats.AveragePrice = 0
	}

	if err := r.db.GetDB().Save(&stats).Error; err != nil {
		return fmt.Errorf("%s: failed to save stats: %w", op, err)
	}

	return nil
}
