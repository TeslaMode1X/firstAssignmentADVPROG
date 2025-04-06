package usecase

import "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"

type ProductUsecase interface {
	ProductDataProcessing(product *model.Product) error
	ProductDataGetAll() ([]model.Product, error)
	ProductDataGetByID(productID int) (*model.Product, error)
	ProductDataUpdate(product *model.Product) error
	ProductDataDelete(productID int) error
}
