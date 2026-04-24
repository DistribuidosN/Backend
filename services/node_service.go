package services

import (
	"archive/zip"
	"Backend/models/interfaces/adapters"
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"Backend/utils"
	"bytes"
	"context"
	"encoding/base64"
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
	resp, err := s.repo.DownloadBatchResult(ctx, token, jobID)
	if err == nil && len(resp.Content) > 0 {
		return resp, nil
	}

	results, resultsErr := s.repo.GetBatchResults(ctx, token, jobID)
	if resultsErr != nil {
		if err != nil {
			return node.BatchDownloadResponse{}, fmt.Errorf("download failed: %w; results fallback failed: %v", err, resultsErr)
		}
		return node.BatchDownloadResponse{}, resultsErr
	}

	if len(results.Images) == 0 {
		if err != nil {
			return node.BatchDownloadResponse{}, fmt.Errorf("download failed and no processed images were returned: %w", err)
		}
		return node.BatchDownloadResponse{}, fmt.Errorf("no processed images available for batch %s", jobID)
	}

	content, zipErr := zipProcessedImages(results.Images)
	if zipErr != nil {
		return node.BatchDownloadResponse{}, zipErr
	}

	return node.BatchDownloadResponse{
		JobID:   results.JobID,
		Content: content,
	}, nil
}

func zipProcessedImages(images []node.ImageItem) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := zip.NewWriter(buffer)

	for _, image := range images {
		name := image.Name
		if name == "" {
			name = fmt.Sprintf("%s.png", image.ID)
		}

		data, err := base64.StdEncoding.DecodeString(image.Base64)
		if err != nil {
			_ = writer.Close()
			return nil, fmt.Errorf("failed to decode processed image %s: %w", name, err)
		}

		entry, err := writer.Create(name)
		if err != nil {
			_ = writer.Close()
			return nil, fmt.Errorf("failed to create zip entry for %s: %w", name, err)
		}

		if _, err := entry.Write(data); err != nil {
			_ = writer.Close()
			return nil, fmt.Errorf("failed to write zip entry for %s: %w", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize zip archive: %w", err)
	}

	return buffer.Bytes(), nil
}
