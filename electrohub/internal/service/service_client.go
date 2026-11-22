package service

import (
	"context"
	"electrohub/internal/domain"
	"electrohub/internal/repository"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type ClientService struct {
	clientRepo  repository.ClientRepository
	addressRepo repository.AddressRepository
}

func NewClientService(c_repo repository.ClientRepository, a_repo repository.AddressRepository) *ClientService {
	return &ClientService{clientRepo: c_repo, addressRepo: a_repo}
}

// ======================================================================================================

func (service *ClientService) CreateClient(ctx context.Context, client domain.Client) (domain.Client, error) {
	if client.Registration_date.IsZero() {
		client.Registration_date = time.Now().UTC()
	}
	address_id, err := service.addressRepo.FindByDetails(ctx, client.Address)
	if err == pgx.ErrNoRows {
		address_id, err = service.addressRepo.Create(ctx, client.Address)
	}
	if err != nil {
		return domain.Client{}, err
	}
	client.Address.ID = address_id
	client_id, err := service.clientRepo.Create(ctx, client)
	if err != nil {
		return domain.Client{}, err
	}
	client.ID = client_id
	return client, nil
}

func (service *ClientService) GetClientsByNameSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	return service.clientRepo.GetByNameSurname(ctx, name, surname)
}

func (service *ClientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return service.clientRepo.Delete(ctx, id)
}

func (service *ClientService) GetClients(ctx context.Context, limit uint, offset uint) ([]domain.Client, error) {
	return service.clientRepo.GetAll(ctx, limit, offset)
}

func (service *ClientService) UpdateClientAddress(ctx context.Context, clientID uuid.UUID, address domain.Address) error {
	address_id, err := service.addressRepo.FindByDetails(ctx, address)
	if err == pgx.ErrNoRows {
		address_id, err = service.addressRepo.Create(ctx, address)
	}
	if err != nil {
		return err
	}
	client, err := service.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}
	address.ID = address_id
	client.Address = address
	return service.clientRepo.Update(ctx, client)
}
