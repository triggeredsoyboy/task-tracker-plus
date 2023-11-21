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

type TaskClient interface {
	CreateTask(token string, task model.Task) (respCode int, err error)
	TaskList(token string) ([]*model.Task, error)
	UpdateTask(token string, id int, task model.UpdateTaskReq) (respCode int, err error)
	DeleteTask(token string, id int) (respCode int, err error)
	TaskByID(token string, id int) (*model.Task, error)
	TaskByCategory(token string) ([]*model.TaskByCategory, error)
}

type taskClient struct {}

func NewTaskClient() *taskClient {
	return &taskClient{}
}

func (t *taskClient) CreateTask(token string, task model.Task) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	dataJson := map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"deadline":    task.Deadline,
		"status":      task.Status,
		"user_id":     task.UserID,
		"category_id": task.CategoryID,
		"priority_id": task.PriorityID,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/task/add"), bytes.NewBuffer(data))
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
		return -1, errors.New("create task failed")
	}

	return resp.StatusCode, nil
}

func (t *taskClient) TaskList(token string) ([]*model.Task, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/task/list"), nil)
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
		return nil, errors.New("cannot get your tasks")
	}

	var tasks []*model.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskClient) UpdateTask(token string, id int, task model.UpdateTaskReq) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	dataJson := map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"deadline":    task.Deadline,
		"status":      task.Status,
		"user_id":     task.UserID,
		"category_id": task.CategoryID,
		"priority_id": task.PriorityID,
	}

	data, err := json.Marshal(dataJson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", config.SetUrl("/api/task/update/"+strconv.Itoa(id)), bytes.NewBuffer(data))
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
		return -1, errors.New("update task failed")
	}

	return resp.StatusCode, nil
}

func (t *taskClient) DeleteTask(token string, id int) (respCode int, err error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("DELETE", config.SetUrl("/api/task/delete/"+strconv.Itoa(id)), nil)
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
		return -1, errors.New("delete task failed")
	}

	return resp.StatusCode, nil
}

func (t *taskClient) TaskByID(token string, id int) (*model.Task, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/task/get/"+strconv.Itoa(id)), nil)
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
		return nil, errors.New("cannot get your task")
	}

	var task *model.Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskClient) TaskByCategory(token string) ([]*model.TaskByCategory, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/task/list-by-category"), nil)
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
		return nil, errors.New("cannot get your tasks")
	}

	var tasks []*model.TaskByCategory
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}