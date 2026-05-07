package services

import (
	"Backend/models/bd"
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/utils"
	"context"
)

type bdService struct {
	repo     adapters.BdRepository
	masterIP string
}

func NewBdService(repo adapters.BdRepository, masterIP string) ports.BdService {
	return &bdService{
		repo:     repo,
		masterIP: masterIP,
	}
}

func (s *bdService) GetPaginatedImages(ctx context.Context, token string, batchUuid string, page int, limit int) (bd.PaginatedImages, error) {
	resp, err := s.repo.GetPaginatedImages(ctx, token, batchUuid, page, limit)
	if err != nil {
		return bd.PaginatedImages{}, err
	}

	for i := range resp.Images {
		resp.Images[i].ResultPath = utils.FixIP(resp.Images[i].ResultPath, s.masterIP)
	}

	return resp, nil
}

func (s *bdService) GetUserBatchesWithCovers(ctx context.Context, token string) ([]bd.BatchWithCover, error) {
	resp, err := s.repo.GetUserBatchesWithCovers(ctx, token)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		resp[i].CoverImageUrl = utils.FixIP(resp[i].CoverImageUrl, s.masterIP)
	}

	return resp, nil
}

func (s *bdService) GetImageMetrics(ctx context.Context, token string, imageUuid string) ([]bd.NodeMetricsDTO, error) {
	return s.repo.GetImageMetrics(ctx, token, imageUuid)
}
