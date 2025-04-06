package repository

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model/dto"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ProductPostgresRepository struct {
	db database.Database
}

func NewProductPostgresRepository(db database.Database) *ProductPostgresRepository {
	return &ProductPostgresRepository{db: db}
}

func (p *ProductPostgresRepository) InsertProduct(product *model.Product) error {
	data := toProductDTO(product)

	result := p.db.GetDb().Create(data)

	if result.Error != nil {
		log.Errorf("InsertProductData: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertProductData: %v rows affected", result.RowsAffected)
	return nil
}

func toProductDTO(prod *model.Product) *dto.ProductDTO {
	var deletedAt gorm.DeletedAt
	if prod.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *prod.DeletedAt, Valid: true}
	} else {
		deletedAt = gorm.DeletedAt{Valid: false}
	}

	return &dto.ProductDTO{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		StockLevel:  prod.StockLevel,
		CategoryID:  prod.CategoryID,
		CreatedAt:   prod.CreatedAt,
		UpdatedAt:   prod.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
