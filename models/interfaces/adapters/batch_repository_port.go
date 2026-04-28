package adapters

import "context"

type BatchRepository interface {
	DownloadBatch(ctx context.Context, token string, batchUuid string) (map[string]interface{}, error)
}
