package repository

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
)

type ProductRepository interface {
	InsertProduct(product *model.Product) error
	GetProducts() ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	ProductUpdate(product *model.Product) error
	ProductDelete(id int) error
	ProductsExists(listOfProducts []string) (bool, error)
}
