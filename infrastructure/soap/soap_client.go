package soap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client es una utilidad genérica para realizar llamadas SOAP (POST text/xml).
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Call ejecuta la petición HTTP POST enviando el body XML y manejando la respuesta.
func (c *Client) Call(url string, xmlBody []byte, token string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	if token != "" {
		authHeader := token
		if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
			authHeader = "Bearer " + token
		}
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en llamada HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	// Permitimos 500 para SOAP Faults
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		return body, fmt.Errorf("error del servidor (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Fault representa un error SOAP (SOAP Fault).
type Fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

func (f *Fault) Reason() string {
	if f.String != "" {
		return f.String
	}
	return "SOAP Fault (unknown error)"
}
