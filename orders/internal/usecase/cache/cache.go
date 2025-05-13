package cache

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
)

type RedisCache interface {
	Get(ctx context.Context, orderID int64) (model.Order, error)
	Set(ctx context.Context, productModel model.Order) error
	Delete(ctx context.Context, productID int64) error
}
