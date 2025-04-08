package usecase

import "github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"

type PromotionUsecase interface {
	ProductDataCreatePromotion(product *model.Promotion) error
	ProductDataGetPromotions() ([]model.Promotion, error)
	ProductDataDeletePromotion(productID string) error
}
