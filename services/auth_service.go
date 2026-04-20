package services

import (
	"Backend/models/auth"
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"context"
)

type authService struct {
	soapRepo adapters.AuthRepository
}

// NewAuthService creates a new instance of the auth service
func NewAuthService(soapRepo adapters.AuthRepository) ports.AuthService {
	return &authService{
		soapRepo: soapRepo,
	}
}

func (s *authService) Login(ctx context.Context, creds auth.UserCredentials) (auth.AuthResponse, error) {
	return s.soapRepo.LogIn(ctx, creds)
}

func (s *authService) Register(ctx context.Context, req auth.RegisterRequest) error {
	return s.soapRepo.Register(ctx, req)
}

func (s *authService) LogOut(ctx context.Context, token string) error {
	return s.soapRepo.LogOut(ctx, token)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (auth.ValidateResponse, error) {
	return s.soapRepo.ValidateToken(ctx, token)
}

func (s *authService) ForgetPwd(ctx context.Context, req auth.ForgetPwdRequest) error {
	return s.soapRepo.ForgetPwd(ctx, req)
}

func (s *authService) ResetPassword(ctx context.Context, token string, req auth.ResetPasswordRequest) error {
	return s.soapRepo.ResetPassword(ctx, token, req)
}
