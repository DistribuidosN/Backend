package ports

import (
	"Backend/models/node"
	"context"
)

// NodeService defines the input port for image processing logic
type NodeService interface {
	UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error)
	UploadBatch(ctx context.Context, token string, req node.BatchUploadRequest) (node.BatchUploadResponse, error)
}
