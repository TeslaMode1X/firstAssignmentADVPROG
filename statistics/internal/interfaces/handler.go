package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

type StatisticsHandler interface {
	GetInventoryStatistics(c *gin.Context)
	GetOrdersStatistics(c *gin.Context)
}

type ProductHandler interface {
	Handler(ctx context.Context, msg *nats.Msg) error
}
