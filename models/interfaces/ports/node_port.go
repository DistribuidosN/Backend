package ports

import (
	"Backend/models/node"
	"context"
	"mime/multipart"
)

// NodeService defines the input port for image processing logic
type NodeService interface {
	UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error)
	ProcessBatch(ctx context.Context, token string, files []*multipart.FileHeader, filters []node.Transformation) (node.BatchJobResponse, error)
	GetBatchStatus(ctx context.Context, token string, jobID string) (node.BatchStatusResponse, error)
	GetBatchResults(ctx context.Context, token string, jobID string) (node.BatchResultsResponse, error)
	GetLogsByImage(ctx context.Context, token string, imageUuid string) ([]node.ProcessingLog, error)
	GetMetricsByNode(ctx context.Context, token string, nodeId string) ([]node.NodeMetric, error)
}
