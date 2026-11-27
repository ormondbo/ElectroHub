package jwtauth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct {
	secretKey []byte
	tokenTTL  time.Duration
}

func NewJwtProvider() (*JwtProvider, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return &JwtProvider{}, errors.New("JWT_SECRET_KEY must be set")
	}
	accessTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		return &JwtProvider{}, errors.New("invalid JWT_ACCESS_TTL")
	}
	return &JwtProvider{
		secretKey: []byte(secret),
		tokenTTL:  accessTTL,
	}, nil
}

// ===================================================================================================================

func (p *JwtProvider) GenerateToken(userID uuid.UUID) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString(p.secretKey)
}

func (p *JwtProvider) ValidateToken(tokenString string) error {
	tokenClaims := &claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		tokenClaims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return p.secretKey, nil
		},
	)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func (p *JwtProvider) GetIDByToken(tokenString string) (uuid.UUID, error) {
	tokenClaims := &claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, tokenClaims)
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(tokenClaims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
