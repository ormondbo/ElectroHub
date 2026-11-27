package postgres

import (
	"time"

	"github.com/google/uuid"
)

type dbUser struct {
	ID           uuid.UUID
	Email        string
	FirstName    string
	LastName     string
	Phone        string
	PasswordHash string
	CreatedAt    time.Time
}
