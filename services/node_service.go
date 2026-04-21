package services

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"Backend/utils"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
)

type nodeService struct {
	repo adapters.NodeRepository
}

func NewNodeService(repo adapters.NodeRepository) ports.NodeService {
	return &nodeService{
		repo: repo,
	}
}

func (s *nodeService) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	return s.repo.UploadImages(ctx, token, req)
}

func (s *nodeService) ProcessBatch(ctx context.Context, token string, files []*multipart.FileHeader, filters []string) (node.BatchJobResponse, error) {
	// 1. Parallel extraction and mapping
	images, err := utils.MapFilesToImageItems(files)
	if err != nil {
		return node.BatchJobResponse{}, fmt.Errorf("error mapping files: %w", err)
	}

	if len(images) == 0 {
		return node.BatchJobResponse{}, fmt.Errorf("no valid images found in the batch")
	}

	// 2. Prepare request with Unique IDs
	req := node.NodeBatchRequest{
		ID:      uuid.NewString(),
		Filters: filters,
		Images:  images,
	}

	// 3. Submit to repository (Now returns a JobID)
	return s.repo.UploadBatch(ctx, token, req)
}

func (s *nodeService) GetBatchStatus(ctx context.Context, token string, jobID string) (node.BatchStatusResponse, error) {
	return s.repo.GetBatchStatus(ctx, token, jobID)
}

func (s *nodeService) GetBatchResults(ctx context.Context, token string, jobID string) (node.BatchResultsResponse, error) {
	// 1. Get results from infrastructure (SOAP)
	results, err := s.repo.GetBatchResults(ctx, token, jobID)
	if err != nil {
		return node.BatchResultsResponse{}, err
	}

	// 2. Apply "Reconversion" Helper (Requested logic)
	reconvertedImages, err := utils.ReconvertProcessedImages(results.Images)
	if err != nil {
		return node.BatchResultsResponse{}, fmt.Errorf("reconversion error: %w", err)
	}

	results.Images = reconvertedImages
	return results, nil
}

func (s *nodeService) DownloadBatchResult(ctx context.Context, token string, jobID string) (node.BatchDownloadResponse, error) {
	return s.repo.DownloadBatchResult(ctx, token, jobID)
}
