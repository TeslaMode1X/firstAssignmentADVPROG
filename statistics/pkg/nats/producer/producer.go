package producer

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/statistics/pkg/nats"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	ProductSubject = "ap2.statistics.event.updated"
	PushTimeout    = time.Second * 30
)

type OrderProducer struct {
	client *nats.Client
}

func NewOrderProducer(client *nats.Client) *OrderProducer {
	return &OrderProducer{
		client: client,
	}
}

func (p *OrderProducer) PublishOrderCreated(ctx context.Context, product model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	err = p.client.Conn.Publish(ProductSubject, data)
	if err != nil {
		return fmt.Errorf("failed to publish product created event: %w", err)
	}

	log.Printf("Published product created event: %d", product.ID)
	return nil
}

func (p *OrderProducer) Push(ctx context.Context, product model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	productPb := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		StockLevel:  int32(product.StockLevel),
		CategoryId:  product.CategoryID,
	}
	data, err := proto.Marshal(productPb)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	err = p.client.Conn.Publish(ProductSubject, data)
	if err != nil {
		return fmt.Errorf("p.product.Conn.Publish: %w", err)
	}
	log.Println("product is pushed:", product)

	return nil
}
