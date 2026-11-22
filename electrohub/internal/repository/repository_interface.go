package repository

import (
	"context"
	"electrohub/internal/domain"
	"errors"
	"uuid"
)

// ======================================================================================================

var ErrNotFound = errors.New("record not found")
var ErrInsufficientStock = errors.New("insufficient product stock")

// ======================================================================================================

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) (uuid.UUID, error)
	Update(ctx context.Context, client domain.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Client, error)
	GetByNameSurname(ctx context.Context, name string, s_name string) ([]domain.Client, error)
	GetAll(ctx context.Context, limit uint, offset uint) ([]domain.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ======================================================================================================

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (uuid.UUID, error)
	Update(ctx context.Context, product domain.Product) error
	DecreaseStock(ctx context.Context, id uuid.UUID, amount uint) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ======================================================================================================

type SupplierRepository interface {
	Create(ctx context.Context, supplier domain.Supplier) (uuid.UUID, error)
	Update(ctx context.Context, supplier domain.Supplier) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error)
	GetAll(ctx context.Context) ([]domain.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ======================================================================================================

type ImageRepository interface {
	Create(ctx context.Context, image domain.Image) (uuid.UUID, error)
	Update(ctx context.Context, image domain.Image) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Image, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ======================================================================================================

type AddressRepository interface {
	Create(ctx context.Context, address domain.Address) (uuid.UUID, error)
	Update(ctx context.Context, address domain.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Address, error)
	FindByDetails(ctx context.Context, address domain.Address) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
