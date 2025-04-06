package repository

import (
	"errors"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model/dto"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
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

func (p *ProductPostgresRepository) GetProducts() ([]model.Product, error) {
	var productDTOs []dto.ProductDTO
	result := p.db.GetDb().Find(&productDTOs)
	if result.Error != nil {
		log.Errorf("GetProducts: %v", result.Error)
		return nil, result.Error
	}

	products := make([]model.Product, len(productDTOs))
	for i, productDTO := range productDTOs {
		products[i] = toProductModel(&productDTO)
	}
	return products, nil
}

func (p *ProductPostgresRepository) GetProductByID(id int) (*model.Product, error) {
	var productDTO dto.ProductDTO
	result := p.db.GetDb().First(&productDTO, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Errorf("GetProductByID: %v", result.Error)
		return nil, result.Error
	}

	product := toProductModel(&productDTO)
	return &product, nil
}

func (p *ProductPostgresRepository) ProductUpdate(product *model.Product) error {
	if product.ID == 0 {
		return fmt.Errorf("product ID cannot be zero")
	}

	updateData := make(map[string]interface{})
	if product.Name != "" {
		updateData["name"] = product.Name
	}
	if product.Description != "" {
		updateData["description"] = product.Description
	}
	if product.Price != 0 {
		updateData["price"] = product.Price
	}
	if product.StockLevel != 0 {
		updateData["stock_level"] = product.StockLevel
	}
	if product.CategoryID != 0 {
		updateData["category_id"] = product.CategoryID
	}
	updateData["updated_at"] = time.Now()

	result := p.db.GetDb().Model(&dto.ProductDTO{}).Where("id = ?", product.ID).Updates(updateData)
	if result.Error != nil {
		log.Errorf("ProductUpdate: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}
	log.Debugf("ProductUpdate: %v rows affected", result.RowsAffected)
	return nil
}

func (p *ProductPostgresRepository) ProductDelete(id int) error {
	if id == 0 {
		return fmt.Errorf("product ID cannot be zero")
	}

	result := p.db.GetDb().Model(&dto.ProductDTO{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		log.Errorf("ProductDelete: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("rows affected: %v", result.RowsAffected)
	}

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

func toProductModel(dto *dto.ProductDTO) model.Product {
	var deletedAt *time.Time
	if dto.DeletedAt.Valid {
		deletedAt = &dto.DeletedAt.Time
	}

	return model.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		StockLevel:  dto.StockLevel,
		CategoryID:  dto.CategoryID,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
