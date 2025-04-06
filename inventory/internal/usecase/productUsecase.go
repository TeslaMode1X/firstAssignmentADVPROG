package usecase

import "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"

type ProductUsecase interface {
	ProductDataProcessing(product *model.Product) error
}
