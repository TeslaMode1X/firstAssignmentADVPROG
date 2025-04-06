package repository

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
)

type ProductRepository interface {
	InsertProduct(product *model.Product) error
}
