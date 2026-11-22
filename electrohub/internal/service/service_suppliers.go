package service

import (
	"context"
	"electrohub/internal/domain"
	"electrohub/internal/repository"

	"uuid"

	"github.com/jackc/pgx/v5"
)

type SupplierService struct {
	supplierRepo repository.SupplierRepository
	addressRepo  repository.AddressRepository
}

func NewSupplierService(supplier_repo repository.SupplierRepository, address_repo repository.AddressRepository) *SupplierService {
	return &SupplierService{supplierRepo: supplier_repo, addressRepo: address_repo}
}

// ======================================================================================================

func (service *SupplierService) CreateSupplier(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error) {
	addressID, err := service.addressRepo.FindByDetails(ctx, supplier.Address)
	if err == pgx.ErrNoRows {
		addressID, err = service.addressRepo.Create(ctx, supplier.Address)
	}
	if err != nil {
		return domain.Supplier{}, err
	}
	supplier.Address.ID = addressID
	supplierID, err := service.supplierRepo.Create(ctx, supplier)
	if err != nil {
		return domain.Supplier{}, err
	}
	supplier.ID = supplierID
	return supplier, nil
}

func (service *SupplierService) UpdateSupplierAddress(ctx context.Context, supplierID uuid.UUID, address domain.Address) error {
	supplier, err := service.supplierRepo.GetByID(ctx, supplierID)
	if err != nil {
		return err
	}
	addressID, err := service.addressRepo.FindByDetails(ctx, address)
	if err == pgx.ErrNoRows {
		addressID, err = service.addressRepo.Create(ctx, address)
	}
	if err != nil {
		return err
	}
	address.ID = addressID
	supplier.Address = address
	return service.supplierRepo.Update(ctx, supplier)

}

func (service *SupplierService) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	return service.supplierRepo.Delete(ctx, id)
}

func (service *SupplierService) GetSuppliers(ctx context.Context) ([]domain.Supplier, error) {
	return service.supplierRepo.GetAll(ctx)
}

func (service *SupplierService) GetSupplierByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error) {
	return service.supplierRepo.GetByID(ctx, id)
}
