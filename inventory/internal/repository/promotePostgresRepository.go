package repository

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model/dto"
	"github.com/labstack/gommon/log"
	"strconv"
)

type PromotePostgresRepository struct {
	db database.Database
}

func NewPromotePostgresRepository(db database.Database) *PromotePostgresRepository {
	return &PromotePostgresRepository{db: db}
}

func (p *PromotePostgresRepository) ProductCreatePromotion(product *model.Promotion) error {
	promotionObj := toPromotionDTO(product)
	result := p.db.GetDb().Create(promotionObj)
	if result.Error != nil {
		log.Errorf("ProductCreatePromotion: %v", result.Error)
		return result.Error
	}
	log.Debugf("ProductCreatePromotion: %v rows affected", result.RowsAffected)

	return nil
}

func (p *PromotePostgresRepository) ProductGetPromotions() ([]model.Promotion, error) {
	var promotionDTOs []dto.PromotionDTO
	result := p.db.GetDb().Find(&promotionDTOs)
	if result.Error != nil {
		log.Errorf("ProductGetPromotions: %v", result.Error)
		return nil, result.Error
	}

	var promotions []model.Promotion
	for _, promotionDTO := range promotionDTOs {
		promotion := toPromotionModel(&promotionDTO)

		productIDs := parseProductIDs(promotionDTO.ApplicableProducts)

		var products []model.Product

		for _, productIDStr := range productIDs {
			productID, err := strconv.Atoi(productIDStr)
			if err != nil {
				log.Errorf("Invalid product ID format: %v", err)
				continue
			}

			var productDTO dto.ProductDTO

			if err := p.db.GetDb().Table("products").Where("id = ?", productID).First(&productDTO).Error; err != nil {
				log.Errorf("Failed to get product with ID %d: %v", productID, err)
				continue
			}

			products = append(products, toProductModel(&productDTO))
		}

		promotion.Products = &products
		promotions = append(promotions, *promotion)
	}

	return promotions, nil
}

func (p *PromotePostgresRepository) ProductDeletePromotion(id string) error {
	if id == "" {
		return fmt.Errorf("promotion ID cannot be empty")
	}

	result := p.db.GetDb().Where("id = ?", id).Delete(&dto.PromotionDTO{})
	if result.Error != nil {
		log.Errorf("ProductDeletePromotion: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("promotion not found with ID: %s", id)
	}

	return nil
}
