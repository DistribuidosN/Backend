package utils

import (
	"Backend/models/node"
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// MapFilesToImageItems handles the in-memory extraction of images from multipart files and archives in parallel.
func MapFilesToImageItems(files []*multipart.FileHeader) ([]node.ImageItem, error) {
	var wg sync.WaitGroup
	resultChan := make(chan []node.ImageItem, len(files))
	errChan := make(chan error, len(files))

	for _, fileHeader := range files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			images, err := extractImagesFromHeader(fh)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- images
		}(fileHeader)
	}

	// Closer goroutine
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	var allImages []node.ImageItem
	for images := range resultChan {
		allImages = append(allImages, images...)
	}

	// For simplicity, we log errors and continue, or we could check errChan.
	for err := range errChan {
		fmt.Printf("Extraction error: %v\n", err)
	}

	return allImages, nil
}

func extractImagesFromHeader(fileHeader *multipart.FileHeader) ([]node.ImageItem, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	switch ext {
	case ".zip":
		return processZip(file, fileHeader.Size)
	case ".tar":
		return processTar(file)
	case ".gz":
		if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".tar.gz") {
			return processTarGz(file)
		}
		return nil, nil
	default:
		if isImageExtension(ext) {
			return processSingleImage(file, fileHeader.Filename)
		}
		return nil, nil
	}
}

func processZip(r io.Reader, size int64) ([]node.ImageItem, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}
	readerAt := bytes.NewReader(buf.Bytes())

	zipReader, err := zip.NewReader(readerAt, int64(buf.Len()))
	if err != nil {
		return nil, err
	}

	var images []node.ImageItem
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() || isIgnoredFile(f.Name) {
			continue
		}

		if isImageExtension(filepath.Ext(f.Name)) {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			img, err := toImageItem(rc, f.Name)
			rc.Close()
			if err == nil {
				images = append(images, img)
			}
		}
	}
	return images, nil
}

func processTar(r io.Reader) ([]node.ImageItem, error) {
	tarReader := tar.NewReader(r)
	var images []node.ImageItem

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return images, err
		}

		if header.Typeflag == tar.TypeReg && !isIgnoredFile(header.Name) && isImageExtension(filepath.Ext(header.Name)) {
			img, err := toImageItem(tarReader, header.Name)
			if err == nil {
				images = append(images, img)
			}
		}
	}
	return images, nil
}

func processTarGz(r io.Reader) ([]node.ImageItem, error) {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	return processTar(gzReader)
}

func processSingleImage(r io.Reader, name string) ([]node.ImageItem, error) {
	if isIgnoredFile(name) {
		return nil, nil
	}
	img, err := toImageItem(r, name)
	if err != nil {
		return nil, err
	}
	return []node.ImageItem{img}, nil
}

func toImageItem(r io.Reader, name string) (node.ImageItem, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return node.ImageItem{}, err
	}

	return node.ImageItem{
		ID:     uuid.NewString(), // Unique ID for each image
		Name:   filepath.Base(name),
		Base64: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

func isImageExtension(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

func isIgnoredFile(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, ".") || strings.Contains(name, "__MACOSX") || strings.EqualFold(base, "thumbs.db")
}
