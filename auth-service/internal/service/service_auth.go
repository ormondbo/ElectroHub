package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/jwtauth"

	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtProvider    *jwtauth.JwtProvider
}

func NewAuthService(userRepository repository.UserRepository, jwtProvider *jwtauth.JwtProvider) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtProvider:    jwtProvider,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	firstName string,
	lastName string,
	phone string,
	password string,
) (domain.User, string, error) {
	_, err := s.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return domain.User{}, "", domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return domain.User{}, "", err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, "", err
	}
	user := domain.User{
		ID:           uuid.New(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return domain.User{}, "", err
	}
	token, err := s.jwtProvider.GenerateToken(user.ID)
	if err != nil {
		return domain.User{}, "", err
	}
	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}
	return s.jwtProvider.GenerateToken(user.ID)
}

func (s *AuthService) ChangePassword(ctx context.Context, tokenString string, oldPassword string, newPassword string) error {
	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return err
	}
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepository.UpdatePassword(ctx, userID, string(passwordHash))
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	if err := s.jwtProvider.ValidateToken(tokenString); err != nil {
		return uuid.Nil, domain.ErrInvalidToken
	}
	userID, err := s.jwtProvider.GetIDByToken(tokenString)
	if err != nil {
		return uuid.Nil, domain.ErrInvalidToken
	}
	return userID, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	randomBytes := make([]byte, 12)
	if _, err := rand.Read(randomBytes); err != nil {
		return err
	}
	newPassword := base64.RawURLEncoding.EncodeToString(randomBytes)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.userRepository.UpdatePassword(ctx, user.ID, string(passwordHash)); err != nil {
		return err
	}
	log.Printf(
		"Password reset for %s. New password: %s",
		email,
		newPassword,
	)
	return nil
}
