package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"first-project/config"
	"first-project/model"
	"io"
	"net/http"
	"strconv"
)

type CategoryClient interface {
	CreateCategory(token string, category model.Category) (respCode int, err error)
	CategoryList(token string) ([]*model.Category, error)
	UpdateCategory(token string, id int, category model.Category) (respCode int, err error)
	DeleteCategory(token string, id int) (respCode int, err error)
	CategoryByID(token string, id int) (*model.Category, error)
}

type categoryClient struct {}

func NewCategoryClient() *categoryClient {
	return &categoryClient{}
}

func (c *categoryClient) CreateCategory(token string, category model.Category) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	dataJson := map[string]interface{}{
		"category_name": category.Name,
		"user_id":       category.UserID,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/category/add"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return -1, errors.New("create category failed")
	}

	return resp.StatusCode, nil
}

func (c *categoryClient) CategoryList(token string) ([]*model.Category, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/category/list"), nil)
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
		return nil, errors.New("cannot get your categories")
	}

	var categories []*model.Category
	err = json.Unmarshal(data, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *categoryClient) UpdateCategory(token string, id int, category model.Category) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	dataJson := map[string]interface{}{
		"category_name": category.Name,
		"user_id":       category.UserID,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", config.SetUrl("/api/category/update/"+strconv.Itoa(id)), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return -1, errors.New("update category failed")
	}

	return resp.StatusCode, nil
}

func (c *categoryClient) DeleteCategory(token string, id int) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("DELETE", config.SetUrl("/api/category/delete/"+strconv.Itoa(id)), nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return -1, errors.New("delete category failed")
	}

	return resp.StatusCode, nil
}

func (c *categoryClient) CategoryByID(token string, id int) (*model.Category, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/category/get/"+strconv.Itoa(id)), nil)
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
		return nil, errors.New("cannot get your category")
	}

	var category *model.Category
	err = json.Unmarshal(data, &category)
	if err != nil {
		return nil, err
	}

	return category, nil
}