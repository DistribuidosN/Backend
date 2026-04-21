package repository

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/node"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type nodeSoapRepository struct {
	url string
}

func NewNodeRepository(url string) adapters.NodeRepository {
	return &nodeSoapRepository{url: url}
}

type nodeSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
	Body    struct {
		UploadImages      *uploadImagesRequest      `xml:"enf:uploadImages,omitempty"`
		UploadImagesBatch *uploadImagesBatchRequest `xml:"enf:uploadImagesBatch,omitempty"`
	} `xml:"soapenv:Body"`
}

type uploadImagesRequest struct {
	ImageData       string                `xml:"imageData"`
	FileName        string                `xml:"fileName"`
	Transformations []node.Transformation `xml:"transformations"`
	Parameters      []node.Transformation `xml:"parameters"`
}

type batchImage struct {
	OriginalName string `xml:"originalName"`
	Data         []byte `xml:"data"`
}

type uploadImagesBatchRequest struct {
	Images          []batchImage          `xml:"images"`
	Transformations []node.Transformation `xml:"transformations"`
	Parameters      []node.Transformation `xml:"parameters"`
}

type nodeSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		UploadImagesResponse      *uploadImagesResponse `xml:"uploadImagesResponse,omitempty"`
		UploadImagesBatchResponse *uploadImagesResponse `xml:"uploadImagesBatchResponse,omitempty"`
	} `xml:"Body"`
}

type uploadImagesResponse struct {
	Return struct {
		ID          int    `xml:"id"`
		UserUUID    string `xml:"userUuid"`
		RequestTime string `xml:"requestTime"`
		StatusID    int    `xml:"statusId"`
	} `xml:"return"`
}

func (r *nodeSoapRepository) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	envelope := &nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	envelope.Body.UploadImages = &uploadImagesRequest{
		ImageData:       req.ImageData,
		FileName:        req.FileName,
		Transformations: req.Transformations,
		Parameters:      []node.Transformation{}, // SOAP expects this field
	}

	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return node.UploadResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(xmlData))
	if err != nil {
		return node.UploadResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	start := time.Now()
	log.Printf("[SOAP][node] --> action=uploadImages url=%s file=%s transformations=%d", r.url, req.FileName, len(req.Transformations))
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[SOAP][node] xx action=uploadImages error=%v", err)
		return node.UploadResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[SOAP][node] <-- action=uploadImages status=%d duration=%s body=%s", resp.StatusCode, time.Since(start).Round(time.Millisecond), string(body))
		return node.UploadResponse{}, fmt.Errorf("soap error: %d", resp.StatusCode)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return node.UploadResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return node.UploadResponse{}, err
	}
	if soapResp.Body.UploadImagesResponse == nil {
		return node.UploadResponse{}, fmt.Errorf("empty uploadImages SOAP response")
	}
	log.Printf("[SOAP][node] <-- action=uploadImages status=%d duration=%s batchId=%d", resp.StatusCode, time.Since(start).Round(time.Millisecond), soapResp.Body.UploadImagesResponse.Return.ID)

	return node.UploadResponse{
		Status:  mapBatchStatus(soapResp.Body.UploadImagesResponse.Return.StatusID),
		Message: fmt.Sprintf("batch created with id %d", soapResp.Body.UploadImagesResponse.Return.ID),
		FileURL: "",
	}, nil
}

func (r *nodeSoapRepository) UploadBatch(ctx context.Context, token string, req node.BatchUploadRequest) (node.BatchUploadResponse, error) {
	images := make([]batchImage, 0, len(req.Files))
	for _, fileHeader := range req.Files {
		image, err := readBatchImage(fileHeader)
		if err != nil {
			return node.BatchUploadResponse{}, err
		}
		images = append(images, image)
	}

	transformations := make([]node.Transformation, 0, len(req.Filters))
	for _, filter := range req.Filters {
		transformations = append(transformations, node.Transformation{Name: filter})
	}

	envelope := &nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	envelope.Body.UploadImagesBatch = &uploadImagesBatchRequest{
		Images:          images,
		Transformations: transformations,
		Parameters:      []node.Transformation{},
	}

	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return node.BatchUploadResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(xmlData))
	if err != nil {
		return node.BatchUploadResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	start := time.Now()
	log.Printf("[SOAP][node] --> action=uploadImagesBatch url=%s files=%d filters=%d", r.url, len(req.Files), len(req.Filters))
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[SOAP][node] xx action=uploadImagesBatch error=%v", err)
		return node.BatchUploadResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[SOAP][node] <-- action=uploadImagesBatch status=%d duration=%s body=%s", resp.StatusCode, time.Since(start).Round(time.Millisecond), string(body))
		return node.BatchUploadResponse{}, fmt.Errorf("soap error: %d", resp.StatusCode)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return node.BatchUploadResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return node.BatchUploadResponse{}, err
	}
	if soapResp.Body.UploadImagesBatchResponse == nil {
		return node.BatchUploadResponse{}, fmt.Errorf("empty uploadImagesBatch SOAP response")
	}
	log.Printf("[SOAP][node] <-- action=uploadImagesBatch status=%d duration=%s batchId=%d", resp.StatusCode, time.Since(start).Round(time.Millisecond), soapResp.Body.UploadImagesBatchResponse.Return.ID)

	return node.BatchUploadResponse{
		Status:  mapBatchStatus(soapResp.Body.UploadImagesBatchResponse.Return.StatusID),
		Message: fmt.Sprintf("batch created with id %d", soapResp.Body.UploadImagesBatchResponse.Return.ID),
	}, nil
}

func readBatchImage(fileHeader *multipart.FileHeader) (batchImage, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return batchImage{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return batchImage{}, err
	}

	return batchImage{
		OriginalName: fileHeader.Filename,
		Data:         data,
	}, nil
}

func mapBatchStatus(statusID int) string {
	switch statusID {
	case 1:
		return "pending"
	case 2:
		return "processing"
	case 3:
		return "finished"
	case 4:
		return "failed"
	default:
		return "accepted"
	}
}
