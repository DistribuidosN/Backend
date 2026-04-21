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
