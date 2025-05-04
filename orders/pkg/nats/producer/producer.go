package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/nats"
	pb "github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	ProductSubject = "orders.order"

	PushTimeout = time.Second * 30
)

type OrderProducer struct {
	client *nats.Client
}

func NewOrdersProducer(client *nats.Client) *OrderProducer {
	return &OrderProducer{
		client: client,
	}
}

func (p *OrderProducer) PublishProductCreated(ctx context.Context, product *pb.Order) error {
	data, err := proto.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal proto product: %w", err)
	}

	err = p.client.Conn.Publish(ProductSubject, data)
	if err != nil {
		return fmt.Errorf("failed to publish proto product: %w", err)
	}

	log.Printf("Published proto product created event: %v", product)
	return nil
}

func (p *OrderProducer) Push(ctx context.Context, product model.Order) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	var productPbItems []*pb.OrderItem

	for _, item := range product.Items {
		productItem := &pb.OrderItem{
			ProductId:   item.ProductID,
			Quantity:    int32(item.Quantity),
			ProductName: item.ProductName,
			Price:       item.Price,
		}
		productPbItems = append(productPbItems, productItem)
	}

	productPb := &pb.Order{
		Id:          0,
		Items:       productPbItems,
		TotalAmount: product.TotalAmount,
		Status:      string(product.Status),
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

func (p *OrderProducer) PublishProductsCreated(ctx context.Context, product []model.Order) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	err = p.client.Conn.Publish(ProductSubject, data)
	if err != nil {
		return fmt.Errorf("failed to publish product created event: %w", err)
	}

	log.Printf("Published product created event: %s", product)
	return nil
}
