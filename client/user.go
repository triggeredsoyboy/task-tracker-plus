package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"first-project/config"
	"first-project/model"
	"io"
	"net/http"
)

type UserClient interface {
	GetCurrentUser(token string) (*model.User, error)
	Register(fullname string, email string, password string) (respCode int, err error)
	Login(email string, password string) (respCode int, err error)
}

type userClient struct{}

func NewUserClient() *userClient {
	return &userClient{}
}

func (u *userClient) GetCurrentUser(token string) (*model.User, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/user/profile"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("cannot get your details")
	}

	var user *model.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userClient) Register(fullname string, email string, password string) (respCode int, err error) {
	dataJson := map[string]string{
		"fullname": fullname,
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/user/register"), bytes.NewBuffer(data))
	if err != nil {
		return -1, nil
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if err != nil {
		return -1, err
	} else {
		return resp.StatusCode, nil
	}
}

func (u *userClient) Login(email string, password string) (respCode int, err error) {
	dataJson := map[string]string{
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/user/login"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if err != nil {
		return -1, err
	} else {
		return resp.StatusCode, nil
	}
}