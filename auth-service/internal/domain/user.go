package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type User struct {
	ID           uuid.UUID
	Email        string
	FirstName    string
	LastName     string
	Phone        string
	PasswordHash string
	CreatedAt    time.Time
}
