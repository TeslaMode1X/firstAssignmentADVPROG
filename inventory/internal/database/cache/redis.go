package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"log"
	"strconv"
	"time"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	keyPrefix = "product:%s"
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

func (c *Client) InitializeCache(ctx context.Context, products []model.Product) error {
	pattern := "product:*"
	// clearing existing cache before initializing
	err := c.clearPattern(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to clear existing cache: %w", err)
	}

	for _, product := range products {
		err := c.Set(ctx, product)
		if err != nil {
			return fmt.Errorf("failed to cache product %d: %w", product.ID, err)
		}
	}

	log.Printf("Successfully initialized cache with %d products", len(products))
	return nil
}

func (c *Client) clearPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.client.Unwrap().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		if len(keys) > 0 {
			if err := c.client.Unwrap().Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}

func (c *Client) Set(ctx context.Context, productModel model.Product) error {
	productModelRedis := model.ProductRedis{
		ID:          productModel.ID,
		Description: productModel.Description,
		CategoryID:  productModel.CategoryID,
		Price:       productModel.Price,
		StockLevel:  productModel.StockLevel,
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
		DeletedAt:   productModel.DeletedAt,
	}

	data, err := json.Marshal(productModelRedis)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	return c.client.Unwrap().Set(ctx, c.key(productModelRedis.ID), data, c.ttl).Err()
}

func (c *Client) Get(ctx context.Context, productID int64) (model.Product, error) {
	data, err := c.client.Unwrap().Get(ctx, c.key(productID)).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return model.Product{}, nil // not found
		}
		return model.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	var productRedis model.ProductRedis
	err = json.Unmarshal(data, &productRedis)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	user := model.Product{
		ID:          productRedis.ID,
		Description: productRedis.Description,
		CategoryID:  productRedis.CategoryID,
		StockLevel:  productRedis.StockLevel,
		Price:       productRedis.Price,
		CreatedAt:   productRedis.CreatedAt,
		UpdatedAt:   productRedis.UpdatedAt,
		DeletedAt:   productRedis.DeletedAt,
	}

	return user, nil
}

func (c *Client) GetAll(ctx context.Context) ([]model.Product, error) {
	var (
		cursor   uint64
		products []model.Product
		pattern  = "product:*"
	)

	for {
		keys, nextCursor, err := c.client.Unwrap().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		for _, key := range keys {
			data, err := c.client.Unwrap().Get(ctx, key).Bytes()
			if err != nil {
				if err == goredis.Nil {
					continue
				}
				return nil, fmt.Errorf("failed to get key %s: %w", key, err)
			}

			var productRedis model.ProductRedis
			if err := json.Unmarshal(data, &productRedis); err != nil {
				return nil, fmt.Errorf("failed to unmarshal product data: %w", err)
			}

			product := model.Product{
				ID:          productRedis.ID,
				Name:        productRedis.Name,
				Description: productRedis.Description,
				Price:       productRedis.Price,
				StockLevel:  productRedis.StockLevel,
				CategoryID:  productRedis.CategoryID,
				CreatedAt:   productRedis.CreatedAt,
				UpdatedAt:   productRedis.UpdatedAt,
				DeletedAt:   productRedis.DeletedAt,
			}

			products = append(products, product)
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return products, nil
}

func (c *Client) Delete(ctx context.Context, productID int64) error {
	return c.client.Unwrap().Del(ctx, c.key(productID)).Err()
}

func (c *Client) key(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	return fmt.Sprintf(keyPrefix, idStr)
}
