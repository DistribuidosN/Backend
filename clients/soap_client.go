package clients

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"Backend/models"
)

// SoapClient is the client to communicate with the SOAP server
type SoapClient struct {
	URL string
}

// NewSoapClient creates a new SoapClient
func NewSoapClient(url string) *SoapClient {
	return &SoapClient{URL: url}
}

// UploadBatch sends the batch of images to the SOAP server
func (c *SoapClient) UploadBatch(req models.BatchRestRequest) (*models.BatchSoapResponse, error) {
	soapReq := models.BatchSoapRequest{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Ser:     "http://service.java.central.com/",
	}
	soapReq.Body.UploadImagesBatch.Token = req.Token
	soapReq.Body.UploadImagesBatch.Images = req.Images
	soapReq.Body.UploadImagesBatch.Compression = req.Compression
	soapReq.Body.UploadImagesBatch.Format = req.Format

	xmlReq, err := xml.Marshal(soapReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(xmlReq))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	httpReq.Header.Set("SOAPAction", "urn:UploadImagesBatch")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var soapResp models.BatchSoapResponse
	err = xml.Unmarshal(body, &soapResp)
	if err != nil {
		return nil, err
	}

	return &soapResp, nil
}
