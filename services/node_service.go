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
	"log"
)

type nodeService struct {
	repo adapters.NodeRepository
}

// NewNodeService creates a new instance of the node service
func NewNodeService(repo adapters.NodeRepository) ports.NodeService {
	return &nodeService{
		repo: repo,
	}
}

func (s *nodeService) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	return s.repo.UploadImages(ctx, token, req)
}

func (s *nodeService) ProcessBatch(ctx context.Context, token string, files []*multipart.FileHeader, filters []node.Transformation) (node.BatchJobResponse, error) {
	log.Printf("[SERVICE] Procesando lote de %d archivos con filtros: %v", len(files), filters)

	// 1. Parallel extraction and mapping
	images, err := utils.MapFilesToImageItems(files)
	if err != nil {
		log.Printf("[SERVICE] Error en la extracción del lote: %v", err)
		return node.BatchJobResponse{}, fmt.Errorf("error mapping files: %w", err)
	}

	log.Printf("[SERVICE] Extracción exitosa. %d imágenes listas para procesar", len(images))

	// 2. Prepare request with Unique IDs
	batchID := uuid.NewString()
	log.Printf("[SERVICE] Generando ID de lote (BatchJobID): %s", batchID)
	req := node.NodeBatchRequest{
		ID:      batchID,
		Filters: filters,
		Images:  images,
	}

	// 3. Submit to repository (Now returns a JobID)
	resp, err := s.repo.UploadBatch(ctx, token, req)
	if err != nil {
		log.Printf("[SERVICE] Error al enviar lote al repositorio SOAP: %v", err)
		return node.BatchJobResponse{}, err
	}

	log.Printf("[SERVICE] Lote enviado con éxito. JobID: %s", resp.BatchId)
	return resp, nil
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

func (s *nodeService) GetLogsByImage(ctx context.Context, token string, imageUuid string) ([]node.ProcessingLog, error) {
	return s.repo.GetLogsByImage(ctx, token, imageUuid)
}

func (s *nodeService) GetMetricsByNode(ctx context.Context, token string, nodeId string) ([]node.NodeMetric, error) {
	return s.repo.GetMetricsByNode(ctx, token, nodeId)
}
