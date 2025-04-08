package usecase

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
	"github.com/labstack/gommon/log"
)

type PromoteUsecaseImpl struct {
	productRepository repository.ProductRepository
	promoteRepository repository.PromoteRepository
}

func NewPromoteUsecaseImpl(productRepository repository.ProductRepository, promoteRepository repository.PromoteRepository) *PromoteUsecaseImpl {
	return &PromoteUsecaseImpl{productRepository: productRepository, promoteRepository: promoteRepository}
}

func (pu *PromoteUsecaseImpl) ProductDataCreatePromotion(product *model.Promotion) error {
	log.Infof("Creating promotion: %+v", product)

	if err := product.Validate(); err != nil {
		log.Errorf("Validation error: %v", err)
		return err
	}

	log.Infof("Checking if products exist: %v", product.ApplicableProducts)
	exists, err := pu.productRepository.ProductsExists(product.ApplicableProducts)
	if err != nil {
		log.Errorf("Error checking products: %v", err)
		return err
	}
	if !exists {
		return fmt.Errorf("products does not exist")
	}

	if err := pu.promoteRepository.ProductCreatePromotion(product); err != nil {
		log.Errorf("Error creating promotion: %v", err)
		return err
	}

	return nil
}

func (pu *PromoteUsecaseImpl) ProductDataGetPromotions() ([]model.Promotion, error) {
	promotionObjs, err := pu.promoteRepository.ProductGetPromotions()
	if err != nil {
		return nil, err
	}

	return promotionObjs, nil
}

func (pu *PromoteUsecaseImpl) ProductDataDeletePromotion(promotionID string) error {
	if promotionID == "" {
		return fmt.Errorf("promotion ID cannot be empty")
	}

	if err := pu.promoteRepository.ProductDeletePromotion(promotionID); err != nil {
		return err
	}
	return nil
}
