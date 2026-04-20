package ports

import (
	"Backend/models/user"
	"context"
)

// UserService defines the input port for user profile and activity logic
type UserService interface {
	GetProfile(ctx context.Context, token string) (user.UserProfile, error)
	UpdateProfile(ctx context.Context, token string, data user.UserProfile) error
	GetActivity(ctx context.Context, token string) ([]user.UserActivity, error)
	SearchUser(ctx context.Context, token string, uid string) (user.UserSearchResponse, error)
	DeleteAccount(ctx context.Context, token string) error
	GetStatistics(ctx context.Context, token string) (user.UserStats, error)
}
