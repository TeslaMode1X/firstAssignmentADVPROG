package repository

import "github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	GetOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	UpdateOrderStatusByID(id int, message string) error
}
