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

// UploadResponse represents the result of an image upload
type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	FileURL string `json:"fileUrl,omitempty"`
}
