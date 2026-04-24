package repository

import (
	"Backend/infrastructure/soap"
	"Backend/models/interfaces/adapters"
	"Backend/models/node"
	"context"
	"encoding/xml"
	"fmt"
	"time"
)

type nodeSoapRepository struct {
	client *soap.Client
	url    string
}

func NewNodeRepository(client *soap.Client, url string) adapters.NodeRepository {
	return &nodeSoapRepository{
		client: client,
		url:    url,
	}
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
		GetLogsByImage  *getLogsByImageRequest  `xml:"enf:getLogsByImage,omitempty"`
		GetMetricsByNode *getMetricsByNodeRequest `xml:"enf:getMetricsByNode,omitempty"`
	} `xml:"soapenv:Body"`
}

type getLogsByImageRequest struct {
	ImageUUID string `xml:"imageUuid"`
}

type getMetricsByNodeRequest struct {
	NodeID string `xml:"nodeId"`
}

type uploadImagesRequest struct {
	ImageData       string                `xml:"imageData"`
	FileName        string                `xml:"fileName"`
	Transformations []node.Transformation `xml:"transformations>transformation"`
}

type uploadBatchWrapper struct {
	Request uploadBatchRequest `xml:"request"`
}

type uploadBatchRequest struct {
	ID      string                `xml:"id"`
	Filters []node.Transformation `xml:"filters"`
	Images  []imageDto            `xml:"images"`
}

type getBatchStatusWrapper struct {
	Request getBatchStatusRequest `xml:"request"`
}

type getBatchStatusRequest struct {
	JobID string `xml:"jobId"`
}

type getBatchResultsWrapper struct {
	Request getBatchResultsRequest `xml:"request"`
}

type getBatchResultsRequest struct {
	JobID string `xml:"jobId"`
}

type imageDto struct {
	ID     string `xml:"id"`
	Name   string `xml:"name"`
	Base64 string `xml:"base64"`
}

// Response Structures
type nodeSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		UploadImagesResponse    *uploadImagesResponse    `xml:"uploadImagesResponse,omitempty"`
		UploadBatchResponse     *uploadBatchResponse     `xml:"uploadBatchResponse,omitempty"`
		GetBatchStatusResponse  *getBatchStatusResponse  `xml:"getBatchStatusResponse,omitempty"`
		GetBatchResultsResponse *getBatchResultsResponse `xml:"getBatchResultsResponse,omitempty"`
		Fault                   *soap.Fault             `xml:"Fault,omitempty"`
	} `xml:"Body"`
}

type uploadImagesResponse struct {
	Return struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
		FileURL string `xml:"fileUrl"`
	} `xml:"return"`
}

type uploadBatchResponse struct {
	Return struct {
		BatchId   string `xml:"batchId"`
		Status  string `xml:"status"`
		Message string `xml:"message"`
	} `xml:"return"`
}

type getBatchStatusResponse struct {
	Return struct {
		BatchId   string `xml:"batchId"`
		Status  string `xml:"status"`
		Message string `xml:"message"`
	} `xml:"return"`
}

type getBatchResultsResponse struct {
	Return struct {
		BatchId  string     `xml:"batchId"`
		Images []imageDto `xml:"images>image"`
	} `xml:"return"`
}

type getLogsByImageResponse struct {
	Return []struct {
		ID        int    `xml:"id"`
		ImageUUID string `xml:"imageUuid"`
		Level     string `xml:"level"`
		Message   string `xml:"message"`
		CreatedAt string `xml:"createdAt"`
	} `xml:"return"`
}

type getMetricsByNodeResponse struct {
	Return []struct {
		ID          int     `xml:"id"`
		NodeID      string  `xml:"nodeId"`
		RAMUsage    float64 `xml:"ramUsage"`
		CPUUsage    float64 `xml:"cpuUsage"`
		BusyWorkers int     `xml:"busyWorkers"`
		ReportedAt  string  `xml:"reportedAt"`
	} `xml:"return"`
}

func (r *nodeSoapRepository) UploadImages(ctx context.Context, token string, req node.ImageUploadRequest) (node.UploadResponse, error) {
	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.UploadImages = &uploadImagesRequest{
		ImageData:       req.ImageData,
		FileName:        req.FileName,
		Transformations: req.Transformations,
	}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return node.UploadResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return node.UploadResponse{}, err
	}

	if soapResp.Body.UploadImagesResponse == nil {
		return node.UploadResponse{}, fmt.Errorf("empty upload response")
	}

	ret := soapResp.Body.UploadImagesResponse.Return
	return node.UploadResponse{
		Status:  ret.Status,
		Message: ret.Message,
		FileURL: ret.FileURL,
	}, nil
}

func (r *nodeSoapRepository) UploadBatch(ctx context.Context, token string, req node.NodeBatchRequest) (node.BatchJobResponse, error) {
	var images []imageDto
	for _, img := range req.Images {
		images = append(images, imageDto{ID: img.ID, Name: img.Name, Base64: img.Base64})
	}

	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.UploadBatch = &uploadBatchWrapper{
		Request: uploadBatchRequest{
			ID:      req.ID,
			Filters: req.Filters,
			Images:  images,
		},
	}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return node.BatchJobResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return node.BatchJobResponse{}, err
	}

	if soapResp.Body.Fault != nil {
		return node.BatchJobResponse{}, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.UploadBatchResponse == nil {
		return node.BatchJobResponse{}, fmt.Errorf("empty batch response")
	}

	ret := soapResp.Body.UploadBatchResponse.Return
	return node.BatchJobResponse{
		BatchId:   ret.BatchId,
		Status:  ret.Status,
		Message: ret.Message,
	}, nil
}

func (r *nodeSoapRepository) GetBatchStatus(ctx context.Context, token string, jobID string) (node.BatchStatusResponse, error) {
	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.GetBatchStatus = &getBatchStatusWrapper{
		Request: getBatchStatusRequest{JobID: jobID},
	}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return node.BatchStatusResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return node.BatchStatusResponse{}, err
	}

	if soapResp.Body.Fault != nil {
		return node.BatchStatusResponse{}, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.GetBatchStatusResponse == nil {
		return node.BatchStatusResponse{}, fmt.Errorf("empty status response")
	}

	ret := soapResp.Body.GetBatchStatusResponse.Return
	return node.BatchStatusResponse{
		BatchId:   ret.BatchId,
		Status:  ret.Status,
		Message: ret.Message,
	}, nil
}

func (r *nodeSoapRepository) GetBatchResults(ctx context.Context, token string, jobID string) (node.BatchResultsResponse, error) {
	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.GetBatchResults = &getBatchResultsWrapper{
		Request: getBatchResultsRequest{JobID: jobID},
	}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return node.BatchResultsResponse{}, err
	}

	var soapResp nodeSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return node.BatchResultsResponse{}, err
	}

	if soapResp.Body.Fault != nil {
		return node.BatchResultsResponse{}, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.GetBatchResultsResponse == nil {
		return node.BatchResultsResponse{}, fmt.Errorf("empty results response")
	}

	ret := soapResp.Body.GetBatchResultsResponse.Return
	var images []node.ImageItem
	for _, img := range ret.Images {
		images = append(images, node.ImageItem{ID: img.ID, Name: img.Name, Base64: img.Base64})
	}

	return node.BatchResultsResponse{
		BatchId:   ret.BatchId,
		Images: images,
	}, nil
}

func (r *nodeSoapRepository) GetLogsByImage(ctx context.Context, token string, imageUuid string) ([]node.ProcessingLog, error) {
	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.GetLogsByImage = &getLogsByImageRequest{ImageUUID: imageUuid}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return nil, err
	}

	var soapResp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			Response *getLogsByImageResponse `xml:"getLogsByImageResponse,omitempty"`
		} `xml:"Body"`
	}
	fmt.Printf("[DEBUG BACKEND] Logs RAW XML: %s\n", string(resp))
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		fmt.Printf("[ERROR BACKEND] Error unmarshalling logs: %v\n", err)
		return nil, err
	}

	if soapResp.Body.Response == nil {
		return nil, fmt.Errorf("empty logs response")
	}

	var logs []node.ProcessingLog
	for _, l := range soapResp.Body.Response.Return {
		created := time.Time{}
		t := parseNodeSOAPTime(l.CreatedAt)
		if t != nil {
			created = *t
		}
		logs = append(logs, node.ProcessingLog{
			ID:        l.ID,
			ImageUUID: l.ImageUUID,
			Level:     l.Level,
			Message:   l.Message,
			CreatedAt: created,
		})
	}
	return logs, nil
}

func (r *nodeSoapRepository) GetMetricsByNode(ctx context.Context, token string, nodeId string) ([]node.NodeMetric, error) {
	env := nodeSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://node.soap.model.server.enfok/",
	}
	env.Body.GetMetricsByNode = &getMetricsByNodeRequest{NodeID: nodeId}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return nil, err
	}

	var soapResp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			Response *getMetricsByNodeResponse `xml:"getMetricsByNodeResponse,omitempty"`
		} `xml:"Body"`
	}
	fmt.Printf("[DEBUG BACKEND] Metrics RAW XML: %s\n", string(resp))
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		fmt.Printf("[ERROR BACKEND] Error unmarshalling metrics: %v\n", err)
		return nil, err
	}

	if soapResp.Body.Response == nil {
		return nil, fmt.Errorf("empty metrics response")
	}

	var metrics []node.NodeMetric
	for _, m := range soapResp.Body.Response.Return {
		reported := time.Time{}
		t := parseNodeSOAPTime(m.ReportedAt)
		if t != nil {
			reported = *t
		}
		metrics = append(metrics, node.NodeMetric{
			ID:          m.ID,
			NodeID:      m.NodeID,
			RAMUsage:    m.RAMUsage,
			CPUUsage:    m.CPUUsage,
			BusyWorkers: m.BusyWorkers,
			ReportedAt:  reported,
		})
	}
	return metrics, nil
}

func parseNodeSOAPTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02T15:04:05",
	}
	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return &t
		}
	}
	return nil
}
