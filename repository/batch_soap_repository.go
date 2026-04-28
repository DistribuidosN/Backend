package repository

import (
	"Backend/infrastructure/soap"
	"context"
	"encoding/xml"
	"fmt"
)

type batchSoapRepository struct {
	client *soap.Client
	url    string
}

func NewBatchSoapRepository(client *soap.Client, url string) *batchSoapRepository {
	return &batchSoapRepository{
		client: client,
		url:    url,
	}
}

// SOAP Structures
type batchSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	BatchNS string   `xml:"xmlns:bat,attr"`
	Body    struct {
		DownloadBatch *downloadBatchRequest `xml:"bat:downloadBatch,omitempty"`
	} `xml:"soapenv:Body"`
}

type downloadBatchRequest struct {
	Request downloadBatchRequestDto `xml:"request"`
}

type downloadBatchRequestDto struct {
	BatchId string `xml:"batchId"`
}

// Response Structures
type batchSoapResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		DownloadBatchResponse *downloadBatchResponseWrapper `xml:"downloadBatchResponse,omitempty"`
		Fault                 *soap.Fault                   `xml:"Fault,omitempty"`
	} `xml:"Body"`
}

type downloadBatchResponseWrapper struct {
	Return downloadBatchResponseDto `xml:"return"`
}

type downloadBatchResponseDto struct {
	DownloadUrl string `xml:"downloadUrl"`
	Status      string `xml:"status"`
	Message     string `xml:"message"`
}

func (r *batchSoapRepository) DownloadBatch(ctx context.Context, token string, batchUuid string) (map[string]interface{}, error) {
	env := batchSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		BatchNS: "http://soap.model.server.enfok/",
	}
	env.Body.DownloadBatch = &downloadBatchRequest{
		Request: downloadBatchRequestDto{
			BatchId: batchUuid,
		},
	}

	xmlData, _ := xml.Marshal(env)
	fmt.Printf("[DEBUG] Request XML (DownloadBatch): %s\n", string(xmlData))

	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		fmt.Printf("[ERROR] SOAP Call failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] Response XML: %s\n", string(resp))

	var soapResp batchSoapResponseEnvelope
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return nil, fmt.Errorf("error al procesar respuesta XML: %w", err)
	}

	if soapResp.Body.Fault != nil {
		return nil, fmt.Errorf("error del orquestador (SOAP Fault): %s", soapResp.Body.Fault.Reason())
	}

	if soapResp.Body.DownloadBatchResponse == nil {
		return nil, fmt.Errorf("respuesta vacía del orquestador")
	}

	return map[string]interface{}{
		"download_url": soapResp.Body.DownloadBatchResponse.Return.DownloadUrl,
		"status":       soapResp.Body.DownloadBatchResponse.Return.Status,
		"message":      soapResp.Body.DownloadBatchResponse.Return.Message,
	}, nil
}
