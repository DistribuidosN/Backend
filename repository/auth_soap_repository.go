package repository

import (
	"Backend/models/auth"
	"Backend/models/interfaces/adapters"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type authSoapRepository struct {
	url string
}

func NewAuthSoapRepository(url string) adapters.AuthRepository {
	return &authSoapRepository{url: url}
}

// SOAP Structures
type soapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
	Header  struct{} `xml:"soapenv:Header"`
	Body    struct {
		LogIn         *logInRequest         `xml:"enf:logIn,omitempty"`
		Register      *registerRequest      `xml:"enf:register,omitempty"`
		LogOut        *struct{}             `xml:"enf:logOut,omitempty"`
		ValidateToken *struct{}             `xml:"enf:validateToken,omitempty"`
		ForgetPwd     *forgetPwdRequest     `xml:"enf:forgetPwd,omitempty"`
		ResetPassword *resetPasswordRequest `xml:"enf:resetPassword,omitempty"`
	} `xml:"soapenv:Body"`
}

type logInRequest struct {
	Email    string `xml:"email"`
	Password string `xml:"password"`
}

type registerRequest struct {
	Email    string `xml:"email"`
	Password string `xml:"password"`
	Username string `xml:"username"`
	RoleId   int    `xml:"role_id"`
}

type forgetPwdRequest struct {
	Email       string `xml:"email"`
	NewPassword string `xml:"newPassword"`
}

type resetPasswordRequest struct {
	NewPassword string `xml:"newPassword"`
}

// Response Structures
type soapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		LogInResponse         *logInResponse         `xml:"logInResponse,omitempty"`
		RegisterResponse      *struct{}              `xml:"registerResponse,omitempty"`
		LogOutResponse        *struct{}              `xml:"logOutResponse,omitempty"`
		ValidateTokenResponse *validateTokenResponse `xml:"validateTokenResponse,omitempty"`
		ForgetPwdResponse     *struct{}              `xml:"forgetPwdResponse,omitempty"`
		ResetPasswordResponse *struct{}              `xml:"resetPasswordResponse,omitempty"`
	} `xml:"Body"`
}

type logInResponse struct {
	Return struct {
		Token  string `xml:"token"`
		RoleId int    `xml:"roleId"` // Java uses CamelCase por defecto en SOAP
	} `xml:"return"`
}

type validateTokenResponse struct {
	Return auth.ValidateResponse `xml:"return"`
}

func (r *authSoapRepository) call(ctx context.Context, token string, bodyFunc func(*soapEnvelope)) (*soapResponse, error) {
	envelope := &soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	bodyFunc(envelope)

	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// Importante: Algunos servidores JAX-WS requieren el header SOAPAction
	req.Header.Set("SOAPAction", "")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		// Imprime el error para debug
		log.Printf("SOAP Error %d: %s", resp.StatusCode, string(respData))
		return nil, fmt.Errorf("soap server error: %d", resp.StatusCode)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var soapResp soapResponse
	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling soap: %w", err)
	}

	return &soapResp, nil
}

func (r *authSoapRepository) LogIn(ctx context.Context, creds auth.UserCredentials) (auth.AuthResponse, error) {
	resp, err := r.call(ctx, "", func(e *soapEnvelope) {
		e.Body.LogIn = &logInRequest{Email: creds.Email, Password: creds.Password}
	})
	if err != nil {
		return auth.AuthResponse{}, err
	}
	return auth.AuthResponse{
		Token:  resp.Body.LogInResponse.Return.Token,
		RoleId: resp.Body.LogInResponse.Return.RoleId,
	}, nil
}

func (r *authSoapRepository) Register(ctx context.Context, req auth.RegisterRequest) error {
	_, err := r.call(ctx, "", func(e *soapEnvelope) {
		e.Body.Register = &registerRequest{
			Email:    req.Email,
			Password: req.Password,
			Username: req.Username,
			RoleId:   req.RoleId,
		}
	})
	return err
}

func (r *authSoapRepository) LogOut(ctx context.Context, token string) error {
	_, err := r.call(ctx, token, func(e *soapEnvelope) {
		e.Body.LogOut = &struct{}{}
	})
	return err
}

func (r *authSoapRepository) ValidateToken(ctx context.Context, token string) (auth.ValidateResponse, error) {
	resp, err := r.call(ctx, token, func(e *soapEnvelope) {
		e.Body.ValidateToken = &struct{}{} // Genera <enf:validateToken/>
	})
	if err != nil {
		return auth.ValidateResponse{}, err
	}

	if resp.Body.ValidateTokenResponse == nil {
		return auth.ValidateResponse{}, fmt.Errorf("empty validate response")
	}

	return resp.Body.ValidateTokenResponse.Return, nil
}

func (r *authSoapRepository) ForgetPwd(ctx context.Context, req auth.ForgetPwdRequest) error {
	_, err := r.call(ctx, "", func(e *soapEnvelope) {
		e.Body.ForgetPwd = &forgetPwdRequest{Email: req.Email, NewPassword: req.NewPassword}
	})
	return err
}

func (r *authSoapRepository) ResetPassword(ctx context.Context, token string, req auth.ResetPasswordRequest) error {
	_, err := r.call(ctx, token, func(e *soapEnvelope) {
		e.Body.ResetPassword = &resetPasswordRequest{NewPassword: req.NewPassword}
	})
	return err
}
