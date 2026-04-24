package node

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
	ID      string      `json:"id"`
	Filters []string    `json:"filters"`
	Images  []ImageItem `json:"images"`
}

// ImageItem represents a single image in base64 format
type ImageItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Base64 string `json:"base64"`
}

// BatchJobResponse represents the initial response after submitting a batch
type BatchJobResponse struct {
	JobID   string `json:"jobId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BatchStatusResponse represents the status of a batch job
type BatchStatusResponse struct {
	JobID   string `json:"jobId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BatchResultsResponse represents the processed images of a batch job
type BatchResultsResponse struct {
	JobID  string      `json:"jobId"`
	Images []ImageItem `json:"images"`
}

type BatchDownloadResponse struct {
	JobID   string
	Content []byte
}

// NodeMetricsResponse - Admin endpoint response for node hardware metrics
type NodeMetricsResponse struct {
	NodeID      string  `json:"nodeId"`
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	GPUUsage    float64 `json:"gpuUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	ActiveJobs  int     `json:"activeJobs"`
	TotalJobs   int     `json:"totalJobs"`
	Uptime      string  `json:"uptime"`
	Timestamp   string  `json:"timestamp"`
}

// ImageLogsResponse - Admin endpoint response for image processing logs
type ImageLogsResponse struct {
	ImageUUID string     `json:"imageUuid"`
	JobID     string     `json:"jobId"`
	Logs      []LogEntry `json:"logs"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}
