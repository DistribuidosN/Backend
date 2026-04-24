package repository

import (
	"Backend/infrastructure/soap"
	"Backend/models/interfaces/adapters"
	"Backend/models/user"
	"context"
	"encoding/xml"
	"fmt"
	"time"
)

type userSoapRepository struct {
	client *soap.Client
	url    string
}

func NewUserSoapRepository(client *soap.Client, url string) adapters.UserRepository {
	return &userSoapRepository{
		client: client,
		url:    url,
	}
}

// SOAP Structures
type userSoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	UserNS  string   `xml:"xmlns:user,attr"`
	Body    struct {
		Profile           *struct{}             `xml:"user:profile,omitempty"`
		UpdateProfile     *updateProfileRequest `xml:"user:updateProfile,omitempty"`
		GetUserActivity   *struct{}             `xml:"user:getUserActivity,omitempty"`
		SearchUser        *searchUserRequest    `xml:"user:searchUser,omitempty"`
		DeleteAccount     *struct{}             `xml:"user:deleteAccount,omitempty"`
		GetUserStatistics *struct{}             `xml:"user:getUserStatistics,omitempty"`
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

// Response Structures
type userSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProfileResponse           *profileResponse           `xml:"profileResponse,omitempty"`
		GetUserActivityResponse   *getUserActivityResponse   `xml:"getUserActivityResponse,omitempty"`
		SearchUserResponse        *searchUserResponse        `xml:"searchUserResponse,omitempty"`
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
		BatchID     string    `xml:"batchId"`
		RequestTime time.Time `xml:"requestTime"`
		Status      string    `xml:"status"`
		ImageCount  int       `xml:"imageCount"`
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
		TotalBatches    int `xml:"totalBatches"`
		TotalImages     int `xml:"totalImages"`
		ImagesCompleted int `xml:"imagesCompleted"`
		ImagesFailed    int `xml:"imagesFailed"`
	} `xml:"return"`
}

func (r *userSoapRepository) GetProfile(ctx context.Context, token string) (user.UserProfile, error) {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.Profile = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return user.UserProfile{}, err
	}

	var soapResp userSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return user.UserProfile{}, err
	}

	if soapResp.Body.ProfileResponse == nil {
		return user.UserProfile{}, fmt.Errorf("empty profile response")
	}

	ret := soapResp.Body.ProfileResponse.Return
	return user.UserProfile{
		ID:        ret.ID,
		UserUUID:  ret.UserUUID,
		Username:  ret.Username,
		Email:     ret.Email,
		RoleID:    ret.RoleID,
		Status:    ret.Status,
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r *userSoapRepository) UpdateProfile(ctx context.Context, token string, data user.UserProfile) error {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.UpdateProfile = &updateProfileRequest{
		UserData: struct {
			Username string `xml:"username"`
			Status   int    `xml:"status"`
		}{Username: data.Username, Status: data.Status},
	}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, token)
	return err
}

func (r *userSoapRepository) GetActivity(ctx context.Context, token string) ([]user.UserActivity, error) {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.GetUserActivity = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return nil, err
	}

	var soapResp userSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return nil, err
	}

	if soapResp.Body.GetUserActivityResponse == nil {
		return nil, fmt.Errorf("empty activity response")
	}

	var activities []user.UserActivity
	for _, a := range soapResp.Body.GetUserActivityResponse.Return {
		activities = append(activities, user.UserActivity{
			BatchID:     a.BatchID,
			RequestTime: a.RequestTime,
			Status:      a.Status,
			ImageCount:  a.ImageCount,
		})
	}
	return activities, nil
}

func (r *userSoapRepository) SearchUser(ctx context.Context, token string, uid string) (user.UserProfile, error) {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.SearchUser = &searchUserRequest{UID: uid}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return user.UserProfile{}, err
	}

	var soapResp userSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return user.UserProfile{}, err
	}

	if soapResp.Body.SearchUserResponse == nil {
		return user.UserProfile{}, fmt.Errorf("user not found")
	}

	ret := soapResp.Body.SearchUserResponse.Return
	return user.UserProfile{
		ID:        ret.ID,
		UserUUID:  ret.UserUUID,
		Username:  ret.Username,
		Email:     ret.Email,
		RoleID:    ret.RoleID,
		Status:    ret.Status,
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r *userSoapRepository) DeleteAccount(ctx context.Context, token string) error {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.DeleteAccount = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	_, err := r.client.Call(r.url, xmlData, token)
	return err
}

func (r *userSoapRepository) GetStatistics(ctx context.Context, token string) (user.UserStatistics, error) {
	env := userSoapEnvelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		UserNS:  "http://user.soap.model.server.enfok/",
	}
	env.Body.GetUserStatistics = &struct{}{}

	xmlData, _ := xml.Marshal(env)
	resp, err := r.client.Call(r.url, xmlData, token)
	if err != nil {
		return user.UserStatistics{}, err
	}

	var soapResp userSoapResponse
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return user.UserStatistics{}, err
	}

	if soapResp.Body.GetUserStatisticsResponse == nil {
		return user.UserStatistics{}, fmt.Errorf("empty statistics response")
	}

	ret := soapResp.Body.GetUserStatisticsResponse.Return
	return user.UserStatistics{
		TotalBatches:    ret.TotalBatches,
		TotalImages:     ret.TotalImages,
		ImagesCompleted: ret.ImagesCompleted,
		ImagesFailed:    ret.ImagesFailed,
	}, nil
}
