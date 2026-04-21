package adapters

import (
	"Backend/models/node"
	"context"
)

// NodeRepository defines the output port for image processing infrastructure
type NodeRepository interface {
	UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error)
	UploadBatch(ctx context.Context, token string, req node.BatchUploadRequest) (node.BatchUploadResponse, error)
}
