package service

import (
	"context"
	"electrohub/internal/api/domain"
	"electrohub/internal/api/repository"
	"errors"
	"fmt"

	"uuid"
)

// ======================================================================================================

var ErrProductAlreadyHasImage = errors.New("product already has an image")

// ======================================================================================================

type ImageService struct {
	imageRepo   repository.ImageRepository
	productRepo repository.ProductRepository
}

func NewImageService(image_repo repository.ImageRepository, product_repo repository.ProductRepository) *ImageService {
	return &ImageService{imageRepo: image_repo, productRepo: product_repo}
}

// ======================================================================================================

func (service *ImageService) AddImage(ctx context.Context, image domain.Image, productID uuid.UUID) error {
	product, err := service.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.Image_id != uuid.Nil() {
		return ErrProductAlreadyHasImage
	}
	img_id, err := service.imageRepo.Create(ctx, image)
	if err != nil {
		return err
	}
	product.Image_id = img_id
	if err := service.productRepo.Update(ctx, product); err != nil {
		if cleanupErr := service.imageRepo.Delete(ctx, img_id); cleanupErr != nil {
			return errors.Join(err, fmt.Errorf("remove unlinked image: %w", cleanupErr))
		}
		return err
	}
	return nil
}

func (service *ImageService) UpdateImage(ctx context.Context, image domain.Image) error {
	return service.imageRepo.Update(ctx, image)
}

func (service *ImageService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return service.imageRepo.Delete(ctx, id)
}

func (service *ImageService) GetImageByID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	return service.imageRepo.GetByID(ctx, id)
}

func (service *ImageService) GetImageByProductID(ctx context.Context, product_id uuid.UUID) (domain.Image, error) {
	product, err := service.productRepo.GetByID(ctx, product_id)
	if err != nil {
		return domain.Image{}, err
	}
	if product.Image_id == uuid.Nil() {
		return domain.Image{}, repository.ErrNotFound
	}
	image, err := service.imageRepo.GetByID(ctx, product.Image_id)
	if err != nil {
		return domain.Image{}, err
	}
	return image, nil
}
