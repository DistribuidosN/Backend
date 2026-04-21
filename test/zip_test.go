package test

import (
	"Backend/utils"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
)

func TestZipExtraction_PruebaZip(t *testing.T) {
	// 1. Path to the ZIP file (relative to the project root, so we use ../)
	zipPath := "../prueba.zip"
	
	// 2. Read the real zip file
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		t.Skipf("Skipping test: could not read %s: %v", zipPath, err)
	}

	// 3. Create a multipart structure in memory to get a valid FileHeader
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Create the "images" part
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filepath.Base(zipPath)))
	h.Set("Content-Type", "application/zip")
	
	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("Failed to create part: %v", err)
	}
	part.Write(zipData)
	writer.Close()

	// 4. Parse it back to get the *multipart.FileHeader
	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10MB limit
	if err != nil {
		t.Fatalf("Failed to parse multipart form: %v", err)
	}

	files := form.File["images"]
	if len(files) == 0 {
		t.Fatalf("No files found in parsed form")
	}

	// 5. Test the mapper logic (which includes parallel processing and ID generation)
	images, err := utils.MapFilesToImageItems(files)
	if err != nil {
		t.Fatalf("Error during mapping: %v", err)
	}

	// 6. Report Results
	t.Logf("Total valid images found in ZIP: %d", len(images))
	
	if len(images) == 0 {
		t.Errorf("Expected at least one image, got 0")
	}

	for i, img := range images {
		t.Logf("[%d] ID: %s | Name: %s | Base64 Size: %d", i+1, img.ID, img.Name, len(img.Base64))
		if img.ID == "" {
			t.Errorf("Image %s has an empty ID", img.Name)
		}
	}
}
