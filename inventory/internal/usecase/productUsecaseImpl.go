package usecase

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/repository"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/usecase/cache"
)

type ProductUsecaseImpl struct {
	productRepository repository.ProductRepository
	promoteRepository repository.PromoteRepository
	redisCacheClient  cache.RedisCache
}

func NewProductUsecaseImpl(productRepository repository.ProductRepository, promoteRepository repository.PromoteRepository, redisCache cache.RedisCache) *ProductUsecaseImpl {
	return &ProductUsecaseImpl{productRepository: productRepository, promoteRepository: promoteRepository, redisCacheClient: redisCache}
}

func (uc *ProductUsecaseImpl) ProductDataProcessing(product *model.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	if err := uc.productRepository.InsertProduct(product); err != nil {
		return err
	}

	if err := uc.redisCacheClient.Set(context.Background(), *product); err != nil {
		return err
	}

	return nil
}

func (uc *ProductUsecaseImpl) ProductDataGetAll() ([]model.Product, error) {
	productsCache, err := uc.redisCacheClient.GetAll(context.Background())
	if err == nil {
		fmt.Println("GOT FROM CACHE")
		return productsCache, err
	}

	fmt.Println("GOT FROM DB")

	products, err := uc.productRepository.GetProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductUsecaseImpl) ProductDataGetByID(productID int) (*model.Product, error) {
	productCache, err := uc.redisCacheClient.Get(context.Background(), int64(productID))
	if err == nil {
		fmt.Println("GOT FROM CACHE")
		return &productCache, err
	}

	fmt.Println("GOT FROM DB")
	product, err := uc.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUsecaseImpl) ProductDataUpdate(product *model.Product) error {
	if product.ID == 0 {
		return fmt.Errorf("product ID cannot be zero")
	}

	if err := uc.productRepository.ProductUpdate(product); err != nil {
		return err
	}

	return nil
}

func (uc *ProductUsecaseImpl) ProductDataDelete(productID int) error {
	if err := uc.productRepository.ProductDelete(productID); err != nil {
		return err
	}

	return nil
}

func (uc *ProductUsecaseImpl) RefreshCache(ctx context.Context) error {
	products, err := uc.productRepository.GetProducts()
	if err != nil {
		return fmt.Errorf("failed to fetch products from repository: %w", err)
	}

	err = uc.redisCacheClient.InitializeCache(ctx, products)
	if err != nil {
		return fmt.Errorf("failed to refresh product cache: %w", err)
	}

	return nil
}
