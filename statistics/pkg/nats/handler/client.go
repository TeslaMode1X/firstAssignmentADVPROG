package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/interfaces"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ProductHandler struct {
	service interfaces.StatisticsService
}

func NewProductHandler(service interfaces.StatisticsService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Handler(ctx context.Context, msg *nats.Msg) error {
	data := msg.Data

	product, err := h.tryParseProtobuf(data)
	if err != nil {
		product, err = h.tryParseJSON(data)
		if err != nil {
			log.Printf("Failed to parse product message: %v", err)
			return fmt.Errorf("failed to parse product message: %w", err)
		}
	}

	err = h.service.RecordProductActivity(product)
	if err != nil {
		log.Printf("Failed to record product activity: %v", err)
		return err
	}

	log.Printf("Successfully processed product activity for ID: %d", product.ID)
	return nil
}

func (h *ProductHandler) tryParseProtobuf(data []byte) (*model.Product, error) {
	var pbProduct pb.Product
	err := proto.Unmarshal(data, &pbProduct)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          pbProduct.Id,
		Name:        pbProduct.Name,
		Description: pbProduct.Description,
		Price:       pbProduct.Price,
		StockLevel:  int(pbProduct.StockLevel),
		CategoryID:  pbProduct.CategoryId,
	}, nil
}

func (h *ProductHandler) tryParseJSON(data []byte) (*model.Product, error) {
	var product model.Product
	err := json.Unmarshal(data, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
