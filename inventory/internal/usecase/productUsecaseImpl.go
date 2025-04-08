package usecase

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
)

type ProductUsecaseImpl struct {
	productRepository repository.ProductRepository
	promoteRepository repository.PromoteRepository
}

func NewProductUsecaseImpl(productRepository repository.ProductRepository, promoteRepository repository.PromoteRepository) *ProductUsecaseImpl {
	return &ProductUsecaseImpl{productRepository: productRepository, promoteRepository: promoteRepository}
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

func (uc *ProductUsecaseImpl) ProductDataGetAll() ([]model.Product, error) {
	products, err := uc.productRepository.GetProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductUsecaseImpl) ProductDataGetByID(productID int) (*model.Product, error) {
	product, err := uc.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUsecaseImpl) ProductDataUpdate(product *model.Product) error {
	if product.ID == 0 {
		return fmt.Errorf("product ID cannot be zero")
	}

	if err := uc.productRepository.ProductUpdate(product); err != nil {
		return err
	}

	return nil
}

func (uc *ProductUsecaseImpl) ProductDataDelete(productID int) error {
	if err := uc.productRepository.ProductDelete(productID); err != nil {
		return err
	}

	return nil
}
