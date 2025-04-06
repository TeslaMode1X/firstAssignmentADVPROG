package usecase

import "github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"

type OrderUsecase interface {
	CreateDataOrder(order *model.Order) error
	GetDataOrders() ([]model.Order, error)
	GetDataOrderByID(id int) (*model.Order, error)
	UpdateDataOrderStatusByID(id int, message string) error
}
