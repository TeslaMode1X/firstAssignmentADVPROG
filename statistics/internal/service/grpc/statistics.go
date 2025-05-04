package grpc

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/statistics"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
)

type StatisticsService struct {
	statistics.UnimplementedStatisticsServiceServer
	statisticsService interfaces.StatisticsService
}

func NewStatisticsService(svc interfaces.StatisticsService) *StatisticsService {
	return &StatisticsService{
		statisticsService: svc,
	}
}

func (s *StatisticsService) GetInventoryStatistics(ctx context.Context, empty *statistics.Empty) (*statistics.GetInventoryStatisticsResponse, error) {
	model, err := s.statisticsService.GetInventoryStatisticsService()
	if err != nil {
		return nil, err
	}

	return &statistics.GetInventoryStatisticsResponse{
		TotalProducts:       int32(model.TotalProducts),
		TotalStock:          int32(model.TotalStock),
		TotalInventoryValue: int32(model.TotalInventoryValue),
	}, nil
}

func (s *StatisticsService) GetOrderStatistics(ctx context.Context, empty *statistics.Empty) (*statistics.GetOrderStatisticsResponse, error) {
	model, err := s.statisticsService.GetOrdersStatisticsService()
	if err != nil {
		return nil, err
	}

	return &statistics.GetOrderStatisticsResponse{
		TotalOrders:  int32(model.TotalOrders),
		TotalRevenue: model.TotalRevenue,
		AveragePrice: model.AveragePrice,
	}, nil
}
