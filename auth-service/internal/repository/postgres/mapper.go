package postgres

import (
	"auth-service/internal/domain"

	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapUserError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	var databaseError *pgconn.PgError
	if errors.As(err, &databaseError) && databaseError.Code == "23505" {
		return domain.ErrEmailAlreadyExists
	}
	return err
}

func (m dbUser) toDomainUser() domain.User {
	return domain.User{
		ID:           m.ID,
		Email:        m.Email,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Phone:        m.Phone,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
	}
}

func toDBUser(user domain.User) dbUser {
	return dbUser{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
}
