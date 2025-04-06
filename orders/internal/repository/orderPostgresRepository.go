package repository

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model/dto"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderPostgresRepository struct {
	db database.Database
}

func NewOrderPostgresRepository(db database.Database) *OrderPostgresRepository {
	return &OrderPostgresRepository{
		db: db,
	}
}

func (r *OrderPostgresRepository) CreateOrder(order *model.Order) error {
	var existingOrder dto.OrderDTO

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	err := r.db.GetDb().Where("total_amount = ? AND created_at > ?",
		order.TotalAmount, fiveMinutesAgo).Order("created_at DESC").First(&existingOrder).Error

	if err == nil {
		var existingItems []dto.OrderItemDTO
		if err := r.db.GetDb().Where("order_id = ?", existingOrder.ID).Find(&existingItems).Error; err == nil {
			if len(existingItems) == len(order.Items) {
				matchCount := 0
				for _, newItem := range order.Items {
					for _, existingItem := range existingItems {
						if newItem.ProductID == existingItem.ProductID &&
							newItem.Quantity == existingItem.Quantity &&
							math.Abs(float64(newItem.Price-existingItem.Price)) < 0.01 {
							matchCount++
							break
						}
					}
				}

				if matchCount == len(order.Items) {
					log.Infof("Detected duplicate order, returning existing order ID: %d", existingOrder.ID)
					order.ID = existingOrder.ID
					return nil
				}
			}
		}
	}

	orderDTO := toOrderDTO(order)

	err = r.db.GetDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(orderDTO).Error; err != nil {
			log.Errorf("CreateOrder: failed to create order: %v", err)
			return err
		}

		order.ID = orderDTO.ID

		for _, item := range order.Items {
			itemDTO := toOrderItemDTO(&item, orderDTO.ID)
			if err := tx.Create(itemDTO).Error; err != nil {
				log.Errorf("CreateOrder: failed to create order item: %v", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Debugf("CreateOrder: order %d created successfully", orderDTO.ID)
	return nil
}

func (r *OrderPostgresRepository) GetOrders() ([]model.Order, error) {
	var orderDTOs []dto.OrderDTO
	result := r.db.GetDb().Preload("Items").Find(&orderDTOs)
	if result.Error != nil {
		log.Errorf("GetOrders: %v", result.Error)
		return nil, result.Error
	}

	orders := make([]model.Order, len(orderDTOs))
	for i, orderDTO := range orderDTOs {
		orders[i] = toOrderModel(&orderDTO)
	}
	return orders, nil
}

func (r *OrderPostgresRepository) GetOrderByID(id int) (*model.Order, error) {
	var orderDTO dto.OrderDTO
	result := r.db.GetDb().Preload("Items").First(&orderDTO, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Errorf("GetOrderByID: %v", result.Error)
		return nil, result.Error
	}

	order := toOrderModel(&orderDTO)
	return &order, nil
}

func (r *OrderPostgresRepository) UpdateOrderStatusByID(id int, status string) error {
	if id == 0 {
		return fmt.Errorf("order ID cannot be zero")
	}

	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	result := r.db.GetDb().Model(&dto.OrderDTO{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		log.Errorf("UpdateOrderStatusByID: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	log.Debugf("UpdateOrderStatusByID: %v rows affected", result.RowsAffected)
	return nil
}

func toOrderDTO(order *model.Order) *dto.OrderDTO {
	return &dto.OrderDTO{
		ID:          order.ID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func toOrderItemDTO(item *model.OrderItem, orderID int64) *dto.OrderItemDTO {
	return &dto.OrderItemDTO{
		OrderID:     orderID,
		ProductID:   item.ProductID,
		Quantity:    item.Quantity,
		Price:       item.Price,
		ProductName: item.ProductName,
	}
}

func toOrderModel(dto *dto.OrderDTO) model.Order {

	items := make([]model.OrderItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		items[i] = model.OrderItem{
			ProductID:   itemDTO.ProductID,
			Quantity:    itemDTO.Quantity,
			Price:       itemDTO.Price,
			ProductName: itemDTO.ProductName,
		}
	}

	return model.Order{
		ID:          dto.ID,
		Status:      model.OrderStatus(dto.Status),
		TotalAmount: dto.TotalAmount,
		Items:       items,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}
