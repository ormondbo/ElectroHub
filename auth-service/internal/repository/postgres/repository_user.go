package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"

	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, first_name, last_name, phone, password_hash, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.Email, user.FirstName, user.LastName,
		user.Phone, user.PasswordHash, user.CreatedAt,
	)
	return mapUserError(err)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var model dbUser

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, first_name, last_name, phone, password_hash, created_at
			FROM users
				WHERE id = $1`,
		id,
	).Scan(
		&model.ID, &model.Email, &model.FirstName, &model.LastName,
		&model.Phone, &model.PasswordHash, &model.CreatedAt,
	)

	if err != nil {
		return domain.User{}, mapUserError(err)
	}

	return model.toDomainUser(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var model dbUser

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, first_name, last_name, phone, password_hash, created_at
			FROM users
				WHERE email = $1`,
		email,
	).Scan(
		&model.ID, &model.Email, &model.FirstName, &model.LastName,
		&model.Phone, &model.PasswordHash, &model.CreatedAt,
	)

	if err != nil {
		return domain.User{}, mapUserError(err)
	}

	return model.toDomainUser(), nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users
			SET password_hash = $1
				WHERE id = $2`,
		passwordHash, id)

	return err
}
