package ports

import (
	"Backend/models/auth"
	"context"
)

// AuthService defines the input port for authentication logic
type AuthService interface {
	Login(ctx context.Context, creds auth.UserCredentials) (auth.AuthResponse, error)
	Register(ctx context.Context, req auth.RegisterRequest) error
	LogOut(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (auth.ValidateResponse, error)
	ForgetPwd(ctx context.Context, req auth.ForgetPwdRequest) error
	ResetPassword(ctx context.Context, token string, req auth.ResetPasswordRequest) error
}
