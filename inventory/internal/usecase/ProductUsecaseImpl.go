package usecase

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
)

type ProductUsecaseImpl struct {
	productRepository repository.ProductRepository
}

func NewProductUsecaseImpl(productRepository repository.ProductRepository) *ProductUsecaseImpl {
	return &ProductUsecaseImpl{productRepository: productRepository}
}

func (uc *ProductUsecaseImpl) ProductDataProcessing(product *model.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	if err := uc.productRepository.InsertProduct(product); err != nil {
		return err
	}

	return nil
}
