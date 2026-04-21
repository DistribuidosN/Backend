package repository

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/node"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type nodeSoapRepository struct {
	url string
}

func NewNodeRepository(url string) adapters.NodeRepository {
	return &nodeSoapRepository{url: url}
}

// SOAP Structures
type nodeSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
	Body    struct {
		UploadImages    *uploadImagesRequest    `xml:"enf:uploadImages,omitempty"` // Lo dejamos igual asumiendo que no falla
		UploadBatch     *uploadBatchWrapper     `xml:"enf:uploadBatch,omitempty"`
		GetBatchStatus  *getBatchStatusWrapper  `xml:"enf:getBatchStatus,omitempty"`
		GetBatchResults *getBatchResultsWrapper `xml:"enf:getBatchResults,omitempty"`
	} `xml:"soapenv:Body"`
}

type uploadBatchRequest struct {
	ID      string           `xml:"id"`
	Filters []string         `xml:"filters"`
	Images  []imageItemBatch `xml:"images"`
}

type getBatchStatusRequest struct {
	JobID string `xml:"jobId"`
}

type getBatchResultsRequest struct {
	JobID string `xml:"jobId"`
}

type imageItemBatch struct {
	ID     string `xml:"id"`
	Name   string `xml:"name"`
	Base64 string `xml:"base64"`
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
		UploadImagesResponse    *uploadImagesResponse    `xml:"uploadImagesResponse,omitempty"`
		UploadBatchResponse     *uploadBatchResponse     `xml:"uploadBatchResponse,omitempty"`
		GetBatchStatusResponse  *getBatchStatusResponse  `xml:"getBatchStatusResponse,omitempty"`
		GetBatchResultsResponse *getBatchResultsResponse `xml:"getBatchResultsResponse,omitempty"`
	} `xml:"Body"`
}

type uploadBatchResponse struct {
	Return struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
		JobID   string `xml:"jobId"`
	} `xml:"return"`
}

type getBatchStatusResponse struct {
	Return struct {
		JobID   string `xml:"jobId"`
		Status  string `xml:"status"`
		Message string `xml:"message"`
	} `xml:"return"`
}

type getBatchResultsResponse struct {
	Return struct {
		JobID  string           `xml:"jobId"`
		Images []imageItemBatch `xml:"images"`
	} `xml:"return"`
}
type uploadImagesResponse struct {
	Return struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
		FileURL string `xml:"fileUrl"`
	} `xml:"return"`
}

// --- NUEVOS WRAPPERS ---
type uploadBatchWrapper struct {
	Request *uploadBatchRequest `xml:"request"`
}

type getBatchStatusWrapper struct {
	Request *getBatchStatusRequest `xml:"request"`
}

type getBatchResultsWrapper struct {
	Request *getBatchResultsRequest `xml:"request"`
}

// Repository Methods

func (r *nodeSoapRepository) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	env := r.newEnvelope()
	env.Body.UploadImages = &uploadImagesRequest{
		ImageData:       req.ImageData,
		FileName:        req.FileName,
		Transformations: req.Transformations,
		Parameters:      []node.Transformation{},
	}
	resp, err := r.doCall(ctx, token, env)
	if err != nil {
		return node.UploadResponse{}, err
	}
	if resp.Body.UploadImagesResponse == nil {
		return node.UploadResponse{}, fmt.Errorf("missing uploadImagesResponse")
	}
	ret := resp.Body.UploadImagesResponse.Return
	return node.UploadResponse{Status: ret.Status, Message: ret.Message, FileURL: ret.FileURL}, nil
}

func (r *nodeSoapRepository) UploadBatch(ctx context.Context, token string, req node.NodeBatchRequest) (node.BatchJobResponse, error) {
	env := r.newEnvelope()
	images := make([]imageItemBatch, len(req.Images))
	for i, img := range req.Images {
		images[i] = imageItemBatch{ID: img.ID, Name: img.Name, Base64: img.Base64}
	}
	env.Body.UploadBatch = &uploadBatchWrapper{
		Request: &uploadBatchRequest{ID: req.ID, Filters: req.Filters, Images: images},
	}
	resp, err := r.doCall(ctx, token, env)
	if err != nil {
		return node.BatchJobResponse{}, err
	}
	if resp.Body.UploadBatchResponse == nil {
		return node.BatchJobResponse{}, fmt.Errorf("missing uploadBatchResponse")
	}
	ret := resp.Body.UploadBatchResponse.Return
	return node.BatchJobResponse{JobID: ret.JobID, Status: ret.Status, Message: ret.Message}, nil
}

func (r *nodeSoapRepository) GetBatchStatus(ctx context.Context, token string, jobID string) (node.BatchStatusResponse, error) {
	env := r.newEnvelope()
	env.Body.GetBatchStatus = &getBatchStatusWrapper{
		Request: &getBatchStatusRequest{JobID: jobID},
	}
	resp, err := r.doCall(ctx, token, env)
	if err != nil {
		return node.BatchStatusResponse{}, err
	}
	if resp.Body.GetBatchStatusResponse == nil {
		return node.BatchStatusResponse{}, fmt.Errorf("missing getBatchStatusResponse")
	}
	ret := resp.Body.GetBatchStatusResponse.Return
	return node.BatchStatusResponse{JobID: ret.JobID, Status: ret.Status, Message: ret.Message}, nil
}

func (r *nodeSoapRepository) GetBatchResults(ctx context.Context, token string, jobID string) (node.BatchResultsResponse, error) {
	env := r.newEnvelope()
	env.Body.GetBatchResults =&getBatchResultsWrapper{
		Request: &getBatchResultsRequest{JobID: jobID},
	}
	resp, err := r.doCall(ctx, token, env)
	if err != nil {
		return node.BatchResultsResponse{}, err
	}
	if resp.Body.GetBatchResultsResponse == nil {
		return node.BatchResultsResponse{}, fmt.Errorf("missing getBatchResultsResponse")
	}
	ret := resp.Body.GetBatchResultsResponse.Return
	images := make([]node.ImageItem, len(ret.Images))
	for i, img := range ret.Images {
		images[i] = node.ImageItem{ID: img.ID, Name: img.Name, Base64: img.Base64}
	}
	return node.BatchResultsResponse{JobID: ret.JobID, Images: images}, nil
}

// Internal Helpers

func (r *nodeSoapRepository) newEnvelope() *nodeSoapEnvelope {
	return &nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
}

func (r *nodeSoapRepository) doCall(ctx context.Context, token string, env *nodeSoapEnvelope) (*nodeSoapResponse, error) {
	xmlData, err := xml.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(data, &soapResp); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &soapResp, nil
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
