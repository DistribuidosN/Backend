package node

import "mime/multipart"

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

type BatchImage struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Base64 string `json:"base64"`
}

type NodeBatchRequest struct {
	ID      string       `json:"id"`
	Filters []string     `json:"filters"`
	Images  []BatchImage `json:"images"`
}

type BatchUploadRequest struct {
	Files   []*multipart.FileHeader
	Filters []string
}

// UploadResponse represents the result of an image upload
type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	FileURL string `json:"fileUrl"`
}

type BatchUploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
