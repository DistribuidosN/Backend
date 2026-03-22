package models

import "encoding/xml"

// BatchRestRequest is the request from the Flutter app
type BatchRestRequest struct {
	Token       string   `json:"token"`
	Images      []string `json:"images"`
	Compression string   `json:"compression"`
	Format      string   `json:"format"`
}

// BatchSoapRequest is the request to the SOAP server
type BatchSoapRequest struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Ser     string   `xml:"xmlns:ser,attr"`
	Header  struct{} `xml:"soapenv:Header"`
	Body    struct {
		UploadImagesBatch struct {
			Token       string   `xml:"ser:token"`
			Images      []string `xml:"ser:images"`
			Compression string   `xml:"ser:compression"`
			Format      string   `xml:"ser:format"`
		} `xml:"ser:UploadImagesBatch"`
	} `xml:"soapenv:Body"`
}

// BatchSoapResponse is the response from the SOAP server
type BatchSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		UploadImagesBatchResponse struct {
			Status      string `xml:"status"`
			ProcessedAt string `xml:"processedAt"`
		} `xml:"UploadImagesBatchResponse"`
	} `xml:"Body"`
}
