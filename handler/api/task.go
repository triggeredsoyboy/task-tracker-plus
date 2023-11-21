package api

import (
	"first-project/model"
	"first-project/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskAPI interface {
	CreateTask(c *gin.Context)
	TaskList(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetByID(c *gin.Context)
	GetByCategory(c *gin.Context)
}

type taskAPI struct {
	userService service.UserService
	categoryService service.CategoryService
	priorityService service.PriorityService
	taskService service.TaskService
}

func NewTaskAPI(userService service.UserService, categoryService service.CategoryService, priorityService service.PriorityService, taskService service.TaskService) *taskAPI {
	return &taskAPI{userService, categoryService, priorityService, taskService}
}

func (a *taskAPI) CreateTask(c *gin.Context) {
	var req model.CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	currentUser := a.userService.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("user not authenticated"))
		return
	}
	req.UserID = currentUser.ID

	category, err := a.categoryService.GetByID(int(req.CategoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("category not found"))
		return
	}
	req.CategoryID = category.ID

	priority, err := a.priorityService.GetByID(int(req.PriorityID))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("priority not found"))
		return
	}
	req.PriorityID = priority.ID

	reqBody := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Status:      req.Status,
		UserID:      req.UserID,
		CategoryID:  req.CategoryID,
		PriorityID:  req.PriorityID,
	}

	err = a.taskService.CreateTask(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created",
		"data": reqBody,
	})
}

func (a *taskAPI) TaskList(c *gin.Context) {
	currentUser := a.userService.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("user not authenticated"))
		return
	}

	tasks, err := a.taskService.TaskList(int(currentUser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (a *taskAPI) UpdateTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid task ID"))
		return
	}

	var updatedTask model.UpdateTaskReq
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	err = a.taskService.UpdateTask(taskID, &updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("update task success"))
}

func (a *taskAPI) DeleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid task ID"))
		return
	}

	err = a.taskService.DeleteTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewErrorResponse("delete task success"))
}

func (a *taskAPI) GetByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid task ID"))
		return
	}

	task, err := a.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, task)
}

func (a *taskAPI) GetByCategory(c *gin.Context) {
	currentUser := a.userService.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("user not authenticated"))
		return
	}

	tasks, err := a.taskService.GetByCategory(int(currentUser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, tasks)
}