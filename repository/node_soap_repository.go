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

type nodeSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
	Body    struct {
		UploadImages *uploadImagesRequest `xml:"enf:uploadImages,omitempty"`
	} `xml:"soapenv:Body"`
}

type uploadImagesRequest struct {
	ImageData       string                `xml:"imageData"`
	FileName        string                `xml:"fileName"`
	Transformations []node.Transformation `xml:"transformations"`
	Parameters      []node.Transformation `xml:"parameters"`
}

type nodeSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		UploadImagesResponse *uploadImagesResponse `xml:"uploadImagesResponse,omitempty"`
	} `xml:"Body"`
}

type uploadImagesResponse struct {
	Return struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
		FileURL string `xml:"fileUrl"`
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
	resp, err := client.Do(httpReq)
	if err != nil {
		return node.UploadResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

	return node.UploadResponse{
		Status:  soapResp.Body.UploadImagesResponse.Return.Status,
		Message: soapResp.Body.UploadImagesResponse.Return.Message,
		FileURL: soapResp.Body.UploadImagesResponse.Return.FileURL,
	}, nil
}
