package services

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/utils"
	"context"
)

type batchService struct {
	repo     adapters.BatchRepository
	masterIP string
}

func NewBatchService(repo adapters.BatchRepository, masterIP string) ports.BatchService {
	return &batchService{
		repo:     repo,
		masterIP: masterIP,
	}
}

func (s *batchService) DownloadBatch(ctx context.Context, token string, batchUuid string) (map[string]interface{}, error) {
	resp, err := s.repo.DownloadBatch(ctx, token, batchUuid)
	if err != nil {
		return nil, err
	}
	return utils.FixIPInObject(resp, s.masterIP).(map[string]interface{}), nil
}
