package repository

import (
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
	result := r.db.GetDB().First(&stats)
	if result.Error != nil {
		return result.Error
	}

	stats.TotalProducts += 1
	stats.TotalStock += p.StockLevel
	stats.TotalInventoryValue += int(p.Price) * p.StockLevel

	result = r.db.GetDB().Save(&stats)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
