package web

import (
	"embed"
	"first-project/client"
	"first-project/model"
	"first-project/service"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskWeb interface {
	TaskPage(c *gin.Context)
	TaskByID(c *gin.Context)
	TaskUpdateProcess(c *gin.Context)
}

type taskWeb struct {
	userClient     client.UserClient
	categoryClient client.CategoryClient
	taskClient     client.TaskClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewTaskWeb(userClient client.UserClient, categoryClient client.CategoryClient, taskClient client.TaskClient, sessionService service.SessionService, embed embed.FS) *taskWeb {
	return &taskWeb{userClient, categoryClient, taskClient, sessionService, embed}
}

func (t *taskWeb) TaskPage(c *gin.Context) {
	status := c.Query("status")
    message := c.Query("message")

	// extracting email from context
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	// get session by email
	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get the current user's details
	userDetail, err := t.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get the current user's categories
	userCategories, err := t.categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get the current user's tasks by category
	userTasks, err := t.taskClient.TaskByCategory(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	var dataTemplate = map[string]interface{}{
		"status":          status,
	    "message":         message,
		"fullname":        userDetail.Fullname,
		"email":           email,
		"user_categories": userCategories,
		"user_tasks":      userTasks,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
		"formatTime": func(t time.Time) string {
			return t.Format(time.DateTime)
		},
		"checkDeadline": func(status string, deadline time.Time) bool {
			today := time.Now()
			if status != "Completed" && today.After(deadline) {
				return true
			}

			return false
		},
	}

	var filepath = path.Join("views", "main", "task.html")
	var header = path.Join("views", "general", "header.html")

	temp, err := template.New("task.html").Funcs(funcMap).ParseFS(t.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
	}
}

func (t *taskWeb) TaskByID(c *gin.Context) {
	status := c.Query("status")
    message := c.Query("message")

	// extracting email from context
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	// get session by email
	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get current user details
	userDetail, err := t.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get current user categories
	userCategories, err := t.categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}
	
	// extracting task id from query parameters
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message=Task%20ID%20is%20missing")
		return
	}

	// convert task id to integer
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// get the current user's task by id
	userTask, err := t.taskClient.TaskByID(session.Token, taskID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	var dataTemplate = map[string]interface{}{
		"status_resp":          status,
	    "message":         message,
		"fullname":        userDetail.Fullname,
		"email":           email,
		"user_categories": userCategories,
		"task_id":         userTask.ID,
		"title":           userTask.Title,
		"description":     userTask.Description,
		"created_at":      userTask.CreatedAt,
		"deadline":        userTask.Deadline,
		"status":          userTask.Status,
		"user_id":         userDetail.ID,
		"category_id":     userTask.CategoryID,
		"priority_id":     userTask.PriorityID,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
	}

	var filepath = path.Join("views", "main", "task-detail.html")
	var header = path.Join("views", "general", "header.html")

	temp, err := template.New("task-detail.html").Funcs(funcMap).ParseFS(t.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
	}
}


func (t *taskWeb) TaskUpdateProcess(c *gin.Context) {
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	userDetail, err := t.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	// extracting task id from query parameters
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message=Task%20ID%20is%20missing")
		return
	}

	// convert task id to integer
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	deadline, err := time.Parse("2006-01-02T15:04", c.Request.FormValue("deadline"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	categoryID, _ := strconv.Atoi(c.Request.FormValue("category"))
	priorityID, _ := strconv.Atoi(c.Request.FormValue("priority"))

	task := model.UpdateTaskReq{
		Title:       c.Request.FormValue("title"),
		Description: c.Request.FormValue("description"),
		Deadline:    deadline,
		Status:      c.Request.FormValue("status"),
		UserID:      userDetail.ID,
		CategoryID:  uint(categoryID),
		PriorityID:  uint(priorityID),
	}

	status, err := t.taskClient.UpdateTask(session.Token, taskID, task)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
		return
	}

	if status == 200 {
		c.Redirect(http.StatusSeeOther, "/client/task?status=success&message=Tugas berhasil diperbarui!")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/task?status=error&message="+err.Error())
	}
}