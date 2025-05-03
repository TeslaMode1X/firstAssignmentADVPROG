package handler

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/handler/response"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatisticsHandler struct {
	statisticsService interfaces.StatisticsService
}

func NewStatisticsHandler(statisticsService interfaces.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHandler) GetInventoryStatistics(c *gin.Context) {
	const op = "handler.statistics.GetInventoryStatistics"

	statistics, err := h.statisticsService.GetInventoryStatisticsService()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, op, err)
		return
	}

	response.Response(c, http.StatusOK, op, statistics)
}

func (h *StatisticsHandler) GetOrdersStatistics(c *gin.Context) {
	const op = "handler.statistics.GetOrdersStatistics"

	statistics, err := h.statisticsService.GetOrdersStatisticsService()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, op, err)
		return
	}

	response.Response(c, http.StatusOK, op, statistics)
}
