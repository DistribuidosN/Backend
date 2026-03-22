package services

import (
	"net/http"

	"Backend/clients"
	"Backend/models"
)

// BatchService is the service for processing batches
type BatchService struct {
	soapClient *clients.SoapClient
}

// NewBatchService creates a new BatchService
func NewBatchService(soapClient *clients.SoapClient) *BatchService {
	return &BatchService{soapClient: soapClient}
}

// ProcessBatch processes the batch of images
func (s *BatchService) ProcessBatch(req models.BatchRestRequest) (*models.BatchSoapResponse, error) {
	if req.Token == "" {
		return nil, models.NewApiError(http.StatusBadRequest, "token is required")
	}
	if len(req.Images) == 0 {
		return nil, models.NewApiError(http.StatusBadRequest, "at least one image is required")
	}

	return s.soapClient.UploadBatch(req)
}
