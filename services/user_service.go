package services

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/models/user"
	"context"
)

type userService struct {
	repo adapters.UserRepository
}

// NewUserService creates a new instance of the User service
func NewUserService(repo adapters.UserRepository) ports.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetProfile(ctx context.Context, token string) (user.UserProfile, error) {
	return s.repo.GetProfile(ctx, token)
}

func (s *userService) UpdateProfile(ctx context.Context, token string, data user.UserProfile) error {
	return s.repo.UpdateProfile(ctx, token, data)
}

func (s *userService) GetActivity(ctx context.Context, token string) ([]user.UserActivity, error) {
	return s.repo.GetActivity(ctx, token)
}

func (s *userService) SearchUser(ctx context.Context, token string, uid string) (user.UserProfile, error) {
	return s.repo.SearchUser(ctx, token, uid)
}

func (s *userService) DeleteAccount(ctx context.Context, token string) error {
	return s.repo.DeleteAccount(ctx, token)
}

func (s *userService) GetStatistics(ctx context.Context, token string) (user.UserStats, error) {
	return s.repo.GetStatistics(ctx, token)
}
