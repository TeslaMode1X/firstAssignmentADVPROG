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

	// Попробуй сначала Order (protobuf -> JSON)
	order, err := h.tryParseOrderProtobuf(data)
	if err == nil {
		log.Printf("ORDER SUCCESS: %+v", order)
		return h.service.RecordOrderActivity(order)
	}

	order, err = h.tryParseOrderJSON(data)
	if err == nil {
		log.Printf("ORDER JSON SUCCESS: %+v", order)
		return h.service.RecordOrderActivity(order)
	}

	// Попробуй Product (protobuf -> JSON)
	product, err := h.tryParseProtobuf(data)
	if err == nil {
		log.Printf("PRODUCT SUCCESS: %+v", product)
		return h.service.RecordProductActivity(product)
	}

	product, err = h.tryParseJSON(data)
	if err == nil {
		log.Printf("PRODUCT JSON SUCCESS: %+v", product)
		return h.service.RecordProductActivity(product)
	}

	log.Printf("Failed to parse message as order or product: %v", err)
	return fmt.Errorf("failed to parse message: %w", err)
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

func (h *ProductHandler) tryParseOrderProtobuf(data []byte) (*model.Order, error) {
	var pbOrder pb.Order
	err := proto.Unmarshal(data, &pbOrder)
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		ID:          int64(pbOrder.Id),
		TotalAmount: pbOrder.TotalAmount,
		Status:      model.OrderStatus(pbOrder.Status),
		Items:       make([]model.OrderItem, 0, len(pbOrder.Items)),
	}

	for _, item := range pbOrder.Items {
		order.Items = append(order.Items, model.OrderItem{
			ProductID:   item.ProductId,
			Quantity:    int(item.Quantity),
			ProductName: item.ProductName,
			Price:       item.Price,
		})
	}

	return order, nil
}

func (h *ProductHandler) tryParseOrderJSON(data []byte) (*model.Order, error) {
	var order model.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
