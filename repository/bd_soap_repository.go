package repository

import (
	"Backend/infrastructure/soap"
	"Backend/models/bd"
	"Backend/models/interfaces/adapters"
	"context"
	"encoding/xml"
	"fmt"
)

type bdSoapRepository struct {
	client *soap.Client
	url    string
}

func NewBdSoapRepository(client *soap.Client, url string) adapters.BdRepository {
	return &bdSoapRepository{
		client: client,
		url:    url,
	}
}

// SOAP Structures
type bdSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	BdNS    string   `xml:"xmlns:bd,attr"`
	Body    struct {
		GetPaginatedImages       *getPaginatedImagesRequest `xml:"bd:getPaginatedImages,omitempty"`
		GetUserBatchesWithCovers *struct{}                  `xml:"bd:getUserBatchesWithCovers,omitempty"`
	} `xml:"soapenv:Body"`
}

type getPaginatedImagesRequest struct {
	BatchUUID string `xml:"batchUuid"`
	Page      int    `xml:"page"`
	Limit     int    `xml:"limit"`
}

// Response Structures
type bdSoapResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		PaginatedImagesResponse   *paginatedResp `xml:"getPaginatedImagesResponse,omitempty"`
		BatchesWithCoversResponse *batchesResp   `xml:"getUserBatchesWithCoversResponse,omitempty"`
		Fault                     *soap.Fault    `xml:"Fault,omitempty"`
	} `xml:"Body"`
}

type paginatedResp struct {
	Return bd.PaginatedImages `xml:"return"`
}

type batchesResp struct {
	Return []bd.BatchWithCover `xml:"return"`
}

func (r *bdSoapRepository) GetPaginatedImages(ctx context.Context, token string, batchUuid string, page int, limit int) (bd.PaginatedImages, error) {
	env := bdSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		BdNS:    "http://bd.soap.model.server.enfok/",
	}
	env.Body.GetPaginatedImages = &getPaginatedImagesRequest{
		BatchUUID: batchUuid,
		Page:      page,
		Limit:     limit,
	}

	xmlData, _ := xml.Marshal(env)
	fmt.Printf("[DEBUG] Request XML (GetPaginatedImages): %s\n", string(xmlData))

	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		fmt.Printf("[ERROR] SOAP Call failed: %v\n", err)
		return bd.PaginatedImages{}, err
	}

	fmt.Printf("[DEBUG] Response XML: %s\n", string(resp))

	var soapResp bdSoapResponseEnvelope
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return bd.PaginatedImages{}, fmt.Errorf("error al procesar respuesta XML: %w", err)
	}

	if soapResp.Body.Fault != nil {
		return bd.PaginatedImages{}, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.PaginatedImagesResponse == nil {
		return bd.PaginatedImages{Images: make([]bd.Image, 0)}, nil
	}

	if soapResp.Body.PaginatedImagesResponse.Return.Images == nil {
		soapResp.Body.PaginatedImagesResponse.Return.Images = make([]bd.Image, 0)
	}

	return soapResp.Body.PaginatedImagesResponse.Return, nil
}

func (r *bdSoapRepository) GetUserBatchesWithCovers(ctx context.Context, token string) ([]bd.BatchWithCover, error) {
	env := bdSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		BdNS:    "http://bd.soap.model.server.enfok/",
	}
	env.Body.GetUserBatchesWithCovers = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	fmt.Printf("[DEBUG] Request XML (GetUserBatchesWithCovers): %s\n", string(xmlData))

	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		fmt.Printf("[ERROR] SOAP Call failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] Response XML: %s\n", string(resp))

	var soapResp bdSoapResponseEnvelope
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return nil, fmt.Errorf("error al procesar respuesta XML: %w", err)
	}

	if soapResp.Body.Fault != nil {
		return nil, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.BatchesWithCoversResponse == nil {
		return make([]bd.BatchWithCover, 0), nil
	}

	if soapResp.Body.BatchesWithCoversResponse.Return == nil {
		return make([]bd.BatchWithCover, 0), nil
	}

	return soapResp.Body.BatchesWithCoversResponse.Return, nil
}
