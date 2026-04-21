package repository

import (
	"Backend/models/interfaces/adapters"
	"Backend/models/user"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type userSoapRepository struct {
	url string
}

func NewUserSoapRepository(url string) adapters.UserRepository {
	return &userSoapRepository{url: url}
}

type userSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Enf     string   `xml:"xmlns:enf,attr"`
	Body    struct {
		Profile           *struct{}             `xml:"enf:profile,omitempty"`
		UpdateProfile     *updateProfileRequest `xml:"enf:updateProfile,omitempty"`
		GetUserActivity   *struct{}             `xml:"enf:getUserActivity,omitempty"`
		SearchUser        *searchUserRequest    `xml:"enf:searchUser,omitempty"`
		DeleteAccount     *struct{}             `xml:"enf:deleteAccount,omitempty"`
		GetUserStatistics *struct{}             `xml:"enf:getUserStatistics,omitempty"`
	} `xml:"soapenv:Body"`
}

type updateProfileRequest struct {
	UserData struct {
		Username string `xml:"username"`
		Status   int    `xml:"status"`
	} `xml:"userData"`
}

type searchUserRequest struct {
	UID string `xml:"uid"`
}

type userSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProfileResponse           *profileResponse           `xml:"profileResponse,omitempty"`
		UpdateProfileResponse     *struct{}                  `xml:"updateProfileResponse,omitempty"`
		GetUserActivityResponse   *getUserActivityResponse   `xml:"getUserActivityResponse,omitempty"`
		SearchUserResponse        *searchUserResponse        `xml:"searchUserResponse,omitempty"`
		DeleteAccountResponse     *struct{}                  `xml:"deleteAccountResponse,omitempty"`
		GetUserStatisticsResponse *getUserStatisticsResponse `xml:"getUserStatisticsResponse,omitempty"`
	} `xml:"Body"`
}

type profileResponse struct {
	Return struct {
		ID        int    `xml:"id"`
		UserUUID  string `xml:"userUuid"`
		Username  string `xml:"username"`
		Email     string `xml:"email"`
		RoleID    int    `xml:"roleId"`
		Status    int    `xml:"status"`
		CreatedAt string `xml:"createdAt"`
	} `xml:"return"`
}

type getUserActivityResponse struct {
	Return []struct {
		ID        string `xml:"id"`
		Action    string `xml:"action"`
		Timestamp string `xml:"timestamp"`
	} `xml:"return"`
}

type searchUserResponse struct {
	Return struct {
		ID        int    `xml:"id"`
		UserUUID  string `xml:"userUuid"`
		Username  string `xml:"username"`
		Email     string `xml:"email"`
		RoleID    int    `xml:"roleId"`
		Status    int    `xml:"status"`
		CreatedAt string `xml:"createdAt"`
	} `xml:"return"`
}

type getUserStatisticsResponse struct {
	Return struct {
		ImagesUploaded int `xml:"imagesUploaded"`
		TotalLogins    int `xml:"totalLogins"`
	} `xml:"return"`
}

func (r *userSoapRepository) call(ctx context.Context, token string, bodyFunc func(*userSoapEnvelope)) (*userSoapResponse, error) {
	envelope := &userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Enf:     "http://user.soap.model.server.enfok/",
	}
	bodyFunc(envelope)
	action := userSOAPAction(envelope)

	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	start := time.Now()
	log.Printf("[SOAP][user] --> action=%s url=%s", action, r.url)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[SOAP][user] xx action=%s error=%v", action, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[SOAP][user] <-- action=%s status=%d duration=%s body=%s", action, resp.StatusCode, time.Since(start).Round(time.Millisecond), string(body))
		return nil, fmt.Errorf("soap error: %d", resp.StatusCode)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var soapResp userSoapResponse
	if err := xml.Unmarshal(respData, &soapResp); err != nil {
		return nil, err
	}
	log.Printf("[SOAP][user] <-- action=%s status=%d duration=%s", action, resp.StatusCode, time.Since(start).Round(time.Millisecond))

	return &soapResp, nil
}

func userSOAPAction(envelope *userSoapEnvelope) string {
	switch {
	case envelope.Body.Profile != nil:
		return "profile"
	case envelope.Body.UpdateProfile != nil:
		return "updateProfile"
	case envelope.Body.GetUserActivity != nil:
		return "getUserActivity"
	case envelope.Body.SearchUser != nil:
		return "searchUser"
	case envelope.Body.DeleteAccount != nil:
		return "deleteAccount"
	case envelope.Body.GetUserStatistics != nil:
		return "getUserStatistics"
	default:
		return "unknown"
	}
}

func (r *userSoapRepository) GetProfile(ctx context.Context, token string) (user.UserProfile, error) {
	resp, err := r.call(ctx, token, func(e *userSoapEnvelope) { e.Body.Profile = &struct{}{} })
	if err != nil {
		return user.UserProfile{}, err
	}
	return user.UserProfile{
		ID:        resp.Body.ProfileResponse.Return.ID,
		UserUUID:  resp.Body.ProfileResponse.Return.UserUUID,
		Username:  resp.Body.ProfileResponse.Return.Username,
		Email:     resp.Body.ProfileResponse.Return.Email,
		RoleID:    resp.Body.ProfileResponse.Return.RoleID,
		Status:    resp.Body.ProfileResponse.Return.Status,
		CreatedAt: resp.Body.ProfileResponse.Return.CreatedAt,
	}, nil
}

func (r *userSoapRepository) UpdateProfile(ctx context.Context, token string, data user.UserProfile) error {
	_, err := r.call(ctx, token, func(e *userSoapEnvelope) {
		e.Body.UpdateProfile = &updateProfileRequest{UserData: struct {
			Username string `xml:"username"`
			Status   int    `xml:"status"`
		}{Username: data.Username, Status: data.Status}}
	})
	return err
}

func (r *userSoapRepository) GetActivity(ctx context.Context, token string) ([]user.UserActivity, error) {
	resp, err := r.call(ctx, token, func(e *userSoapEnvelope) { e.Body.GetUserActivity = &struct{}{} })
	if err != nil {
		return nil, err
	}
	activities := make([]user.UserActivity, len(resp.Body.GetUserActivityResponse.Return))
	for i, a := range resp.Body.GetUserActivityResponse.Return {
		activities[i] = user.UserActivity{ID: a.ID, Action: a.Action, Timestamp: a.Timestamp}
	}
	return activities, nil
}

func (r *userSoapRepository) SearchUser(ctx context.Context, token string, uid string) (user.UserProfile, error) {
	resp, err := r.call(ctx, token, func(e *userSoapEnvelope) { e.Body.SearchUser = &searchUserRequest{UID: uid} })
	if err != nil {
		return user.UserProfile{}, err
	}
	return user.UserProfile{
		ID:        resp.Body.SearchUserResponse.Return.ID,
		UserUUID:  resp.Body.SearchUserResponse.Return.UserUUID,
		Username:  resp.Body.SearchUserResponse.Return.Username,
		Email:     resp.Body.SearchUserResponse.Return.Email,
		RoleID:    resp.Body.SearchUserResponse.Return.RoleID,
		Status:    resp.Body.SearchUserResponse.Return.Status,
		CreatedAt: resp.Body.SearchUserResponse.Return.CreatedAt,
	}, nil
}

func (r *userSoapRepository) DeleteAccount(ctx context.Context, token string) error {
	_, err := r.call(ctx, token, func(e *userSoapEnvelope) { e.Body.DeleteAccount = &struct{}{} })
	return err
}

func (r *userSoapRepository) GetStatistics(ctx context.Context, token string) (user.UserStats, error) {
	resp, err := r.call(ctx, token, func(e *userSoapEnvelope) { e.Body.GetUserStatistics = &struct{}{} })
	if err != nil {
		return user.UserStats{}, err
	}
	return user.UserStats{ImagesUploaded: resp.Body.GetUserStatisticsResponse.Return.ImagesUploaded, TotalLogins: resp.Body.GetUserStatisticsResponse.Return.TotalLogins}, nil
}
