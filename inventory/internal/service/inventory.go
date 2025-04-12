package service

import (
	"context"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/proto/gen/inventory"
	"strconv"
	"strings"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryService struct {
	inventory.UnimplementedInventoryServiceServer
	productUsecase usecase.ProductUsecase
}

func NewInventoryService(productUsecase usecase.ProductUsecase) *InventoryService {
	return &InventoryService{
		productUsecase: productUsecase,
	}
}

func (s *InventoryService) CreateProduct(ctx context.Context, req *inventory.CreateProductRequest) (*inventory.ProductResponse, error) {
	categoryID, err := strconv.ParseUint(req.Category, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid category ID: %v", err)
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float32(req.Price),
		StockLevel:  int(req.Stock),
		CategoryID:  int64(uint(categoryID)),
	}

	if err := s.productUsecase.ProductDataProcessing(product); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create product: %v", err)
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
			Stock:       int32(product.StockLevel),
			Category:    req.Category,
		},
	}, nil
}

func (s *InventoryService) GetProductByID(ctx context.Context, req *inventory.GetProductRequest) (*inventory.ProductResponse, error) {
	id, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID: %v", err)
	}

	product, err := s.productUsecase.ProductDataGetByID(int(id))
	if err != nil {
		// Проверяем, что продукт действительно не найден
		if strings.Contains(err.Error(), "record not found") {
			return nil, status.Errorf(codes.NotFound, "Product with ID %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get product: %v", err)
	}

	// Дополнительная проверка на nil
	if product == nil {
		return nil, status.Errorf(codes.NotFound, "Product with ID %s not found", req.Id)
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
			Stock:       int32(product.StockLevel),
			Category:    strconv.FormatUint(uint64(product.CategoryID), 10),
		},
	}, nil
}

func (s *InventoryService) UpdateProduct(ctx context.Context, req *inventory.UpdateProductRequest) (*inventory.ProductResponse, error) {
	id, err := strconv.ParseUint(strconv.FormatInt(req.Id, 10), 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID: %v", err)
	}

	categoryID, err := strconv.ParseUint(req.Category, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid category ID: %v", err)
	}

	product := &model.Product{
		ID:          int64(uint(id)),
		Name:        req.Name,
		Description: req.Description,
		Price:       float32(req.Price),
		StockLevel:  int(req.Stock),
		CategoryID:  int64(uint(categoryID)),
	}

	if err := s.productUsecase.ProductDataUpdate(product); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update product: %v", err)
	}

	return &inventory.ProductResponse{
		Product: &inventory.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
			Stock:       int32(product.StockLevel),
			Category:    req.Category,
		},
	}, nil
}

func (s *InventoryService) DeleteProduct(ctx context.Context, req *inventory.DeleteProductRequest) (*inventory.Empty, error) {
	id, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID: %v", err)
	}

	if err := s.productUsecase.ProductDataDelete(int(id)); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete product: %v", err)
	}

	return &inventory.Empty{}, nil
}

func (s *InventoryService) GetProducts(ctx context.Context, _ *inventory.Empty) (*inventory.GetProductsResponse, error) {
	products, err := s.productUsecase.ProductDataGetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get products: %v", err)
	}

	protoProducts := make([]*inventory.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       float64(p.Price),
			Stock:       int32(p.StockLevel),
			Category:    strconv.FormatUint(uint64(p.CategoryID), 10),
		})
	}

	return &inventory.GetProductsResponse{
		Products: protoProducts,
	}, nil
}
