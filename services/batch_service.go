package services

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"context"
)

type batchService struct {
	repo adapters.BatchRepository
}

func NewBatchService(repo adapters.BatchRepository) ports.BatchService {
	return &batchService{repo: repo}
}

func (s *batchService) DownloadBatch(ctx context.Context, token string, batchUuid string) (map[string]interface{}, error) {
	return s.repo.DownloadBatch(ctx, token, batchUuid)
}
