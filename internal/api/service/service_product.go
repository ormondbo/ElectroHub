package service

import (
	"context"
	"electrohub/internal/api/domain"
	"electrohub/internal/api/repository"

	"uuid"
)

// ======================================================================================================

var ErrInsufficientStock = repository.ErrInsufficientStock

// ======================================================================================================

type ProductService struct {
	productRepo repository.ProductRepository
	imageRepo   repository.ImageRepository
}

func NewProductService(repo repository.ProductRepository, imageRepo repository.ImageRepository) *ProductService {
	return &ProductService{productRepo: repo, imageRepo: imageRepo}
}

// ======================================================================================================

func (service *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if product.Image_id != uuid.Nil() {
		if _, err := service.imageRepo.GetByID(ctx, product.Image_id); err != nil {
			return domain.Product{}, err
		}
	}
	productID, err := service.productRepo.Create(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}
	product.ID = productID
	return product, nil
}

func (service *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return service.productRepo.GetByID(ctx, id)
}

func (service *ProductService) GetAvailableProducts(ctx context.Context) ([]domain.Product, error) {
	return service.productRepo.GetAll(ctx)
}

func (service *ProductService) DecreaseProductStock(ctx context.Context, id uuid.UUID, amount uint) error {
	return service.productRepo.DecreaseStock(ctx, id, amount)
}

func (service *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return service.productRepo.Delete(ctx, id)
}
