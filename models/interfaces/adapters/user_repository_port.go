package adapters

import (
	"Backend/models/user"
	"context"
)

// UserRepository defines the output port for user infrastructure
type UserRepository interface {
	GetProfile(ctx context.Context, token string) (user.UserProfile, error)
	UpdateProfile(ctx context.Context, token string, data user.UserProfile) error
	GetActivity(ctx context.Context, token string) ([]user.UserActivity, error)
	SearchUser(ctx context.Context, token string, uid string) (user.UserProfile, error)
	DeleteAccount(ctx context.Context, token string) error
	GetStatistics(ctx context.Context, token string) (user.UserStats, error)
}
