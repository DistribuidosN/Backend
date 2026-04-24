package node

import "time"


// ImageUploadRequest represents the data required to upload and transform images
type ImageUploadRequest struct {
	ImageData       string           `json:"imageData"`
	FileName        string           `json:"fileName"`
	Transformations []Transformation `json:"transformations"`
}

// Transformation represents a single image transformation
type Transformation struct {
	Name   string `json:"name" xml:"name"`
	Params string `json:"params,omitempty" xml:"params,omitempty"`
}

// UploadResponse represents the result of a single image upload
type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	FileURL string `json:"fileUrl,omitempty"`
}

// NodeBatchRequest represents a batch of images and filters for processing
type NodeBatchRequest struct {
	ID      string           `json:"id"`
	Filters []Transformation `json:"filters"`
	Images  []ImageItem      `json:"images"`
}

// ImageItem represents a single image in base64 format
type ImageItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Base64 string `json:"base64"`
}

// BatchJobResponse represents the initial response after submitting a batch
type BatchJobResponse struct {
	BatchId   string `json:"batchId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BatchStatusResponse represents the status of a batch job
type BatchStatusResponse struct {
	BatchId   string `json:"batchId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BatchResultsResponse represents the processed images of a batch job
type BatchResultsResponse struct {
	BatchId  string      `json:"batchId"`
	Images []ImageItem `json:"images"`
}
// NodeMetric representa la lectura de hardware de un nodo
type NodeMetric struct {
	ID           int       `json:"id"`
	NodeID       string    `json:"node_id"`
	RAMUsage     float64   `json:"ram_usage"`
	CPUUsage     float64   `json:"cpu_usage"`
	BusyWorkers  int       `json:"busy_workers"`
	ReportedAt   time.Time `json:"reported_at"`
}

// ProcessingLog representa un evento de log individual asociado a una imagen
type ProcessingLog struct {
	ID        int       `json:"id"`
	ImageUUID string    `json:"image_uuid"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
