package service

import (
	"context"

	authpb "electrohub/internal/proto/auth"
)

type AuthService struct {
	client authpb.AuthServiceClient
}

func NewAuthService(client authpb.AuthServiceClient) *AuthService {
	return &AuthService{client: client}
}

// ======================================================================================================

func (service *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return service.client.Register(ctx, req)
}
func (service *AuthService) Authenticate(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return service.client.Login(ctx, req)
}
func (service *AuthService) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	return service.client.ResetPassword(ctx, req)
}
func (service *AuthService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	return service.client.ValidateToken(ctx, req)
}
func (service *AuthService) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	return service.client.ChangePassword(ctx, req)
}
