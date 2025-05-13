package usecase

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/handler/client"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase/cache"
	"time"
)

type OrderUsecaseImpl struct {
	orderRepository  repository.OrderRepository
	inventoryClient  client.InventoryClient
	redisCacheClient cache.RedisCache
}

func NewOrderUsecaseImpl(repo repository.OrderRepository, invClient client.InventoryClient, redisCache cache.RedisCache) OrderUsecase {
	return &OrderUsecaseImpl{
		orderRepository:  repo,
		inventoryClient:  invClient,
		redisCacheClient: redisCache,
	}
}

func (uc *OrderUsecaseImpl) CreateDataOrder(order *model.Order) error {
	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	var totalAmount float32 = 0

	for i, item := range order.Items {
		product, err := uc.inventoryClient.GetProductByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product info: %w", err)
		}

		if product.StockLevel < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s (ID: %d): have %d, need %d",
				product.Name, product.ID, product.StockLevel, item.Quantity)
		}

		order.Items[i].Price = product.Price
		order.Items[i].ProductName = product.Name

		totalAmount += product.Price * float32(item.Quantity)
	}

	order.TotalAmount = totalAmount

	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if err := uc.orderRepository.CreateOrder(order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := uc.redisCacheClient.Set(context.Background(), *order); err != nil {
		return fmt.Errorf("failed to set order: %w", err)
	}

	for _, item := range order.Items {
		if err := uc.inventoryClient.UpdateProductStock(item.ProductID, item.Quantity); err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}
	}

	return nil
}

func (uc *OrderUsecaseImpl) GetDataOrders() ([]model.Order, error) {
	orders, err := uc.orderRepository.GetOrders()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (uc *OrderUsecaseImpl) GetDataOrderByID(id int) (*model.Order, error) {
	orderCache, err := uc.redisCacheClient.Get(context.Background(), int64(id))
	if err == nil {
		fmt.Println("GOT FROM REDIS CACHE")
		return &orderCache, nil
	}

	fmt.Println("GOT FROM DATABASE")
	order, err := uc.orderRepository.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *OrderUsecaseImpl) UpdateDataOrderStatusByID(id int, message string) error {
	if err := uc.orderRepository.UpdateOrderStatusByID(id, message); err != nil {
		return err
	}

	return nil
}
