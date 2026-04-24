package adapters

import (
	"Backend/models/bd"
	"context"
)

// BdRepository define el puerto de salida para el almacenamiento y consulta de imágenes.
type BdRepository interface {
	GetPaginatedImages(ctx context.Context, token string, batchUuid string, page int, limit int) (bd.PaginatedImages, error)
	GetUserBatchesWithCovers(ctx context.Context, token string) ([]bd.BatchWithCover, error)
}
