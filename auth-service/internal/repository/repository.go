package repository

import (
	"auth-service/internal/domain"

	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}
