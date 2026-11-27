package api

import (
	"auth-service/internal/domain"
	"auth-service/internal/service"
	authpb "contracts/proto/auth"

	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (s *Handler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetFirstName() == "" || req.GetLastName() == "" || req.GetPhone() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email, first name, last name, phone and password are required")
	}
	user, token, err := s.authService.Register(
		ctx, req.GetEmail(),
		req.GetFirstName(),
		req.GetLastName(),
		req.GetPhone(),
		req.GetPassword(),
	)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &authpb.RegisterResponse{
		UserId:      user.ID.String(),
		Message:     "user registered successfully",
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

func (s *Handler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to authenticate")
	}
	return &authpb.LoginResponse{AccessToken: token, TokenType: "Bearer"}, nil
}

func (s *Handler) ValidateToken(_ context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	if req.GetAccessToken() == "" {
		return &authpb.ValidateTokenResponse{Valid: false}, nil
	}
	userID, err := s.authService.ValidateToken(req.GetAccessToken())
	if err != nil {
		return &authpb.ValidateTokenResponse{Valid: false}, nil
	}
	return &authpb.ValidateTokenResponse{Valid: true, UserId: userID.String()}, nil
}

func (s *Handler) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	if req.GetAccessToken() == "" ||
		req.GetOldPassword() == "" ||
		req.GetNewPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "access token, old password and new password are required")
	}
	err := s.authService.ChangePassword(ctx, req.GetAccessToken(), req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid access token")
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid old password")
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to change password")
	}
	return &authpb.ChangePasswordResponse{Message: "password changed successfully"}, nil
}

func (s *Handler) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	err := s.authService.ResetPassword(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to reset password")
	}
	return &authpb.ResetPasswordResponse{Message: "password reset successfully"}, nil
}
