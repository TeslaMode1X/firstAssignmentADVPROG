package service

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/promotion"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PromotionService struct {
	promotion.UnimplementedPromotionServiceServer
	promotionUsecase usecase.PromotionUsecase
}

func NewPromotionService(promotionUsecase usecase.PromotionUsecase) *PromotionService {
	return &PromotionService{
		promotionUsecase: promotionUsecase,
	}
}

func (s *PromotionService) CreatePromotion(ctx context.Context, req *promotion.CreatePromotionRequest) (*promotion.PromotionResponse, error) {
	if req.Name == "" || req.DiscountPercentage <= 0 || len(req.ApplicableProducts) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid promotion data: name, discount, and applicable products are required")
	}

	promo := &model.Promotion{
		ID:                 req.Id,
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
		ApplicableProducts: req.ApplicableProducts,
		StartDate:          req.StartDate.AsTime(),
		EndDate:            req.EndDate.AsTime(),
		IsActive:           true,
	}

	if err := promo.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Promotion validation failed: %v", err)
	}

	if err := s.promotionUsecase.ProductDataCreatePromotion(promo); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create promotion: %v", err)
	}

	return &promotion.PromotionResponse{
		Promotion: &promotion.Promotion{
			Id:                 promo.ID,
			Name:               promo.Name,
			Description:        promo.Description,
			DiscountPercentage: promo.DiscountPercentage,
			ApplicableProducts: promo.ApplicableProducts,
			Products:           nil,
			StartDate:          timestamppb.New(promo.StartDate),
			EndDate:            timestamppb.New(promo.EndDate),
		},
	}, nil
}

func (s *PromotionService) GetPromotions(ctx context.Context, _ *promotion.Empty) (*promotion.GetPromotionsResponse, error) {
	promotions, err := s.promotionUsecase.ProductDataGetPromotions()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get promotions: %v", err)
	}

	protoPromotions := make([]*promotion.Promotion, 0, len(promotions))
	for _, promo := range promotions {
		var protoProducts []*promotion.Product
		if promo.Products != nil {
			for _, p := range *promo.Products {
				protoProducts = append(protoProducts, &promotion.Product{
					Id:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       float64(p.Price),
					Stock:       int32(p.StockLevel),
					Category:    string(p.CategoryID),
				})
			}
		}

		protoPromotions = append(protoPromotions, &promotion.Promotion{
			Id:                 promo.ID,
			Name:               promo.Name,
			Description:        promo.Description,
			DiscountPercentage: promo.DiscountPercentage,
			ApplicableProducts: promo.ApplicableProducts,
			Products:           protoProducts,
			StartDate:          timestamppb.New(promo.StartDate),
			EndDate:            timestamppb.New(promo.EndDate),
		})
	}

	return &promotion.GetPromotionsResponse{
		Promotions: protoPromotions,
	}, nil
}

func (s *PromotionService) DeletePromotion(ctx context.Context, req *promotion.DeletePromotionRequest) (*promotion.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Promotion ID cannot be empty")
	}

	if err := s.promotionUsecase.ProductDataDeletePromotion(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete promotion: %v", err)
	}

	return &promotion.Empty{}, nil
}
