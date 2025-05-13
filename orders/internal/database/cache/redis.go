package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"strconv"
	"time"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	keyPrefix = "order:%s"
)

type Client struct {
	client *redis.Client
	ttl    time.Duration
}

func NewClient(client *redis.Client, ttl time.Duration) *Client {
	return &Client{
		client: client,
		ttl:    ttl,
	}
}

func (c *Client) Set(ctx context.Context, productModel model.Order) error {
	orderModelRedis := model.OrderRedis{
		ID:          productModel.ID,
		Items:       productModel.Items,
		TotalAmount: productModel.TotalAmount,
		Status:      productModel.Status,
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
		CompletedAt: productModel.CompletedAt,
	}

	data, err := json.Marshal(orderModelRedis)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	return c.client.Unwrap().Set(ctx, c.key(orderModelRedis.ID), data, c.ttl).Err()
}

func (c *Client) Get(ctx context.Context, orderID int64) (model.Order, error) {
	data, err := c.client.Unwrap().Get(ctx, c.key(orderID)).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return model.Order{}, nil // not found
		}
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var orderRedis model.OrderRedis
	err = json.Unmarshal(data, &orderRedis)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	order := model.Order{
		ID:          orderRedis.ID,
		Items:       orderRedis.Items,
		TotalAmount: orderRedis.TotalAmount,
		Status:      orderRedis.Status,
		CreatedAt:   orderRedis.CreatedAt,
		UpdatedAt:   orderRedis.UpdatedAt,
		CompletedAt: orderRedis.CompletedAt,
	}

	return order, nil
}

//func (c *Client) GetAll(ctx context.Context) ([]model.Product, error) {
//	var (
//		cursor   uint64
//		products []model.Product
//		pattern  = "product:*"
//	)
//
//	for {
//		keys, nextCursor, err := c.client.Unwrap().Scan(ctx, cursor, pattern, 100).Result()
//		if err != nil {
//			return nil, fmt.Errorf("failed to scan keys: %w", err)
//		}
//
//		for _, key := range keys {
//			data, err := c.client.Unwrap().Get(ctx, key).Bytes()
//			if err != nil {
//				if err == goredis.Nil {
//					continue
//				}
//				return nil, fmt.Errorf("failed to get key %s: %w", key, err)
//			}
//
//			var productRedis model.ProductRedis
//			if err := json.Unmarshal(data, &productRedis); err != nil {
//				return nil, fmt.Errorf("failed to unmarshal product data: %w", err)
//			}
//
//			product := model.Product{
//				ID:          productRedis.ID,
//				Name:        productRedis.Name,
//				Description: productRedis.Description,
//				Price:       productRedis.Price,
//				StockLevel:  productRedis.StockLevel,
//				CategoryID:  productRedis.CategoryID,
//				CreatedAt:   productRedis.CreatedAt,
//				UpdatedAt:   productRedis.UpdatedAt,
//				DeletedAt:   productRedis.DeletedAt,
//			}
//
//			products = append(products, product)
//		}
//
//		if nextCursor == 0 {
//			break
//		}
//		cursor = nextCursor
//	}
//
//	return products, nil
//}

func (c *Client) Delete(ctx context.Context, orderID int64) error {
	return c.client.Unwrap().Del(ctx, c.key(orderID)).Err()
}

func (c *Client) key(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	return fmt.Sprintf(keyPrefix, idStr)
}
