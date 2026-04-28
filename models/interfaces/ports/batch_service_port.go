package ports

import "context"

type BatchService interface {
	DownloadBatch(ctx context.Context, token string, batchUuid string) (map[string]interface{}, error)
}
