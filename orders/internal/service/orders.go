package service

import (
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/internal/usecase"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/pkg/nats/producer"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/orders"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersService struct {
	orders.UnimplementedOrderServiceServer
	ordersUsecase usecase.OrderUsecase
	producer      *producer.OrderProducer
}

func NewOrdersService(orderUsecase usecase.OrderUsecase, producer *producer.OrderProducer) *OrdersService {
	return &OrdersService{
		ordersUsecase: orderUsecase,
		producer:      producer,
	}
}

func (s *OrdersService) CreateOrder(ctx context.Context, orderProto *orders.Order) (*orders.CreateOrderResponse, error) {
	modelOrder := &model.Order{
		Items: make([]model.OrderItem, 0, len(orderProto.Items)),
	}

	for _, item := range orderProto.Items {
		modelOrder.Items = append(modelOrder.Items, model.OrderItem{
			ProductID:   item.ProductId,
			Quantity:    int(item.Quantity),
			ProductName: item.ProductName,
			Price:       item.Price,
		})
	}

	if err := s.ordersUsecase.CreateDataOrder(modelOrder); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	responseOrder := &orders.Order{
		Status: string(modelOrder.Status),
		Items:  make([]*orders.OrderItem, 0, len(modelOrder.Items)),
	}

	for _, item := range modelOrder.Items {
		responseOrder.Items = append(responseOrder.Items, &orders.OrderItem{
			ProductId:   item.ProductID,
			Quantity:    int32(item.Quantity),
			ProductName: item.ProductName,
			Price:       item.Price,
		})
	}

	if err := s.producer.PublishProductCreated(ctx, responseOrder); err != nil {
		// Логируем ошибку, но не прерываем выполнение
	}

	return &orders.CreateOrderResponse{
		Order: responseOrder,
	}, nil
}

func (s *OrdersService) GetOrder(ctx context.Context, empty *orders.Empty) (*orders.GetOrderResponse, error) {
	modelOrders, err := s.ordersUsecase.GetDataOrders()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get data: %v", err)
	}
	protoOrders := make([]*orders.Order, 0, len(modelOrders))

	for _, modelOrder := range modelOrders {
		protoItems := make([]*orders.OrderItem, 0, len(modelOrder.Items))

		for _, item := range modelOrder.Items {
			protoItems = append(protoItems, &orders.OrderItem{
				ProductId:   item.ProductID,
				Quantity:    int32(item.Quantity),
				ProductName: item.ProductName,
				Price:       item.Price,
			})
		}

		protoOrders = append(protoOrders, &orders.Order{
			Items:  protoItems,
			Status: string(modelOrder.Status),
		})
	}

	if err := s.producer.PublishProductsCreated(ctx, modelOrders); err != nil {
		// Логируем ошибку, но не прерываем выполнение
	}

	return &orders.GetOrderResponse{
		Order: protoOrders,
	}, nil
}

func (s *OrdersService) GetOrderById(ctx context.Context, req *orders.GetOrderRequest) (*orders.GetOrderByIdResponse, error) {
	modelOrder, err := s.ordersUsecase.GetDataOrderByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order by ID: %v", err)
	}

	protoItems := make([]*orders.OrderItem, 0, len(modelOrder.Items))

	for _, item := range modelOrder.Items {
		protoItems = append(protoItems, &orders.OrderItem{
			ProductId:   item.ProductID,
			Quantity:    int32(item.Quantity),
			ProductName: item.ProductName,
			Price:       item.Price,
		})
	}

	protoOrder := &orders.Order{
		Items:  protoItems,
		Status: string(modelOrder.Status),
	}

	return &orders.GetOrderByIdResponse{
		Orders: protoOrder,
	}, nil
}

func (s *OrdersService) UpdateOrderStatusById(ctx context.Context, req *orders.UpdateOrderStatusByIdRequest) (*orders.Empty, error) {
	err := s.ordersUsecase.UpdateDataOrderStatusByID(int(req.Id), req.Message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &orders.Empty{}, nil
}
