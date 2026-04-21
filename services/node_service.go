package services

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"context"
)

type nodeService struct {
	repo adapters.NodeRepository
}

// NewNodeService creates a new instance of the Node service
func NewNodeService(repo adapters.NodeRepository) ports.NodeService {
	return &nodeService{
		repo: repo,
	}
}

func (s *nodeService) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	return s.repo.UploadImages(ctx, token, req)
}

func (s *nodeService) UploadBatch(ctx context.Context, token string, req node.BatchUploadRequest) (node.BatchUploadResponse, error) {
	return s.repo.UploadBatch(ctx, token, req)
}
