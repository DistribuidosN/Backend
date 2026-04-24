package repository

import (
	"Backend/infrastructure/soap"
	"Backend/models/auth"
	"Backend/models/interfaces/adapters"
	"context"
	"encoding/xml"
	"fmt"
)

type authSoapRepository struct {
	client *soap.Client
	url    string
}

func NewAuthSoapRepository(client *soap.Client, url string) adapters.AuthRepository {
	return &authSoapRepository{
		client: client,
		url:    url,
	}
}

// SOAP Structures
type soapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
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
		RoleId int    `xml:"roleId"`
	} `xml:"return"`
}

type validateTokenResponse struct {
	Return auth.ValidateResponse `xml:"return"`
}

func (r *authSoapRepository) LogIn(ctx context.Context, creds auth.UserCredentials) (auth.AuthResponse, error) {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.LogIn = &logInRequest{Email: creds.Email, Password: creds.Password}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, "")
	if err != nil {
		return auth.AuthResponse{}, err
	}

	var soapResp soapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return auth.AuthResponse{}, err
	}

	if soapResp.Body.LogInResponse == nil {
		return auth.AuthResponse{}, fmt.Errorf("login failed")
	}

	ret := soapResp.Body.LogInResponse.Return
	return auth.AuthResponse{
		Token:  ret.Token,
		RoleId: ret.RoleId,
	}, nil
}

func (r *authSoapRepository) Register(ctx context.Context, req auth.RegisterRequest) error {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.Register = &registerRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		RoleId:   req.RoleId,
	}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, "")
	return err
}

func (r *authSoapRepository) LogOut(ctx context.Context, token string) error {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.LogOut = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, token)
	return err
}

func (r *authSoapRepository) ValidateToken(ctx context.Context, token string) (auth.ValidateResponse, error) {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.ValidateToken = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return auth.ValidateResponse{}, err
	}

	var soapResp soapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return auth.ValidateResponse{}, err
	}

	if soapResp.Body.ValidateTokenResponse == nil {
		return auth.ValidateResponse{}, fmt.Errorf("invalid token")
	}

	return soapResp.Body.ValidateTokenResponse.Return, nil
}

func (r *authSoapRepository) ForgetPwd(ctx context.Context, req auth.ForgetPwdRequest) error {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.ForgetPwd = &forgetPwdRequest{Email: req.Email, NewPassword: req.NewPassword}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, "")
	return err
}

func (r *authSoapRepository) ResetPassword(ctx context.Context, token string, req auth.ResetPasswordRequest) error {
	env := soapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://auth.soap.model.server.enfok/",
	}
	env.Body.ResetPassword = &resetPasswordRequest{NewPassword: req.NewPassword}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, token)
	return err
}
