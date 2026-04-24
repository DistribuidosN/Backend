package ports

import (
	"Backend/models/node"
	"context"
	"mime/multipart"
)

// NodeService defines the input port for image processing logic
type NodeService interface {
	UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error)
	ProcessBatch(ctx context.Context, token string, files []*multipart.FileHeader, filters []string) (node.BatchJobResponse, error)
	GetBatchStatus(ctx context.Context, token string, jobID string) (node.BatchStatusResponse, error)
	GetBatchResults(ctx context.Context, token string, jobID string) (node.BatchResultsResponse, error)
	DownloadBatchResult(ctx context.Context, token string, jobID string) (node.BatchDownloadResponse, error)
	
	// Admin endpoints
	GetNodeMetrics(ctx context.Context, token string, nodeID string) (node.NodeMetricsResponse, error)
	GetImageProcessingLogs(ctx context.Context, token string, imageUUID string) (node.ImageLogsResponse, error)
}
