package cache

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
)

type RedisCache interface {
	Get(ctx context.Context, productID int64) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Set(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, productID int64) error
	InitializeCache(ctx context.Context, products []model.Product) error
}
