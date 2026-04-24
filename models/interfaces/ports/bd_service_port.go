package ports

import (
	"Backend/models/bd"
	"context"
)

// BdService define los puertos de entrada para la gestión de galería y lotes.
type BdService interface {
	GetPaginatedImages(ctx context.Context, token string, batchUuid string, page int, limit int) (bd.PaginatedImages, error)
	GetUserBatchesWithCovers(ctx context.Context, token string) ([]bd.BatchWithCover, error)
	GetImageMetrics(ctx context.Context, token string, imageUuid string) ([]bd.NodeMetricsDTO, error)
}
