package service

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
)

type StatisticsService struct {
	statisticsRepo interfaces.StatisticsRepo
}

func NewStatisticsService(statisticsRepo interfaces.StatisticsRepo) *StatisticsService {
	return &StatisticsService{
		statisticsRepo: statisticsRepo,
	}
}

func (s *StatisticsService) GetInventoryStatisticsService() (*model.InventoryStatistics, error) {
	const op = "service.statistics.GetInventoryStatisticsService"

	statistics, err := s.statisticsRepo.GetInventoryStatisticsRepo()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return statistics, nil
}

func (s *StatisticsService) GetOrdersStatisticsService() (*model.OrderStatistics, error) {
	const op = "service.statistics.GetOrdersStatisticsService"

	statistics, err := s.statisticsRepo.GetOrdersStatisticsRepo()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return statistics, nil
}

func (s *StatisticsService) RecordProductActivity(p *model.Product) error {
	const op = "service.statistics.RecordProductActivity"

	err := s.statisticsRepo.RecordProductActivityRepo(p)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
