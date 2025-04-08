package repository

import "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"

type PromoteRepository interface {
	ProductCreatePromotion(product *model.Promotion) error
	ProductGetPromotions() ([]model.Promotion, error)
	ProductDeletePromotion(id string) error
}
