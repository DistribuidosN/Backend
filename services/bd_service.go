package services

import (
	"Backend/models/bd"
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"context"
)

type bdService struct {
	repo adapters.BdRepository
}

func NewBdService(repo adapters.BdRepository) ports.BdService {
	return &bdService{repo: repo}
}

func (s *bdService) GetPaginatedImages(ctx context.Context, token string, batchUuid string, page int, limit int) (bd.PaginatedImages, error) {
	return s.repo.GetPaginatedImages(ctx, token, batchUuid, page, limit)
}

func (s *bdService) GetUserBatchesWithCovers(ctx context.Context, token string) ([]bd.BatchWithCover, error) {
	return s.repo.GetUserBatchesWithCovers(ctx, token)
}
