package utils

import (
	"Backend/models/node"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

// ReconvertProcessedImages handles backend-side transformation and persistence of processed images.
func ReconvertProcessedImages(images []node.ImageItem) ([]node.ImageItem, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images to reconvert")
	}

	// 1. Ensure the output directory exists
	outDir := "images"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create images directory: %w", err)
	}

	// 2. Process and save each image
	for _, img := range images {
		// Decode Base64
		data, err := base64.StdEncoding.DecodeString(img.Base64)
		if err != nil {
			fmt.Printf("Error decoding image %s: %v\n", img.Name, err)
			continue
		}

		// Save to disk
		filePath := filepath.Join(outDir, img.Name)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			fmt.Printf("Error saving image %s to disk: %v\n", img.Name, err)
			continue
		}
		
		fmt.Printf("Successfully saved processed image: %s\n", filePath)
	}

	return images, nil
}
