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

type DashboardWeb interface {
	Dashboard(c *gin.Context)
	TaskAddProcess(c *gin.Context)
	TaskDeleteProcess(c *gin.Context)
}

type dashboardWeb struct {
	userClient     client.UserClient
	categoryClient client.CategoryClient
	taskClient     client.TaskClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewDashboardWeb(userClient client.UserClient, categoryClient client.CategoryClient, taskClient client.TaskClient, sessionService service.SessionService, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{userClient, categoryClient, taskClient, sessionService, embed}
}

func (d *dashboardWeb) Dashboard(c *gin.Context) {
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
	session, err := d.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	// get the current user's details
	userDetail, err := d.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	// get the current user's categories
	userCategories, err := d.categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	// get the current user's tasks by category
	userTasks, err := d.taskClient.TaskByCategory(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
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
		"getOnGoingTasks": func(status string) int {
			var count = 0
			for _, task := range userTasks {
				if task.Status == "On Going" {
					count++
				}
			}
			return count
		},
		"getCompletedTasks": func(status string) int {
			var count = 0
			for _, task := range userTasks {
				if task.Status == "Completed" {
					count++
				}
			}
			return count
		},
	}

	var filepath = path.Join("views", "main", "dashboard.html")
	var header = path.Join("views", "general", "header.html")

	temp, err := template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
	}
}

func (d *dashboardWeb) TaskAddProcess(c *gin.Context) {
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	session, err := d.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	userDetail, err := d.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	deadline, err := time.Parse("2006-01-02T15:04", c.Request.FormValue("deadline"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	categoryID, _ := strconv.Atoi(c.Request.FormValue("category"))
	priorityID, _ := strconv.Atoi(c.Request.FormValue("priority"))

	task := model.Task{
		Title:       c.Request.FormValue("title"),
		Description: c.Request.FormValue("description"),
		Deadline:    deadline,
		Status:      c.Request.FormValue("status"),
		UserID:      userDetail.ID,
		CategoryID:  uint(categoryID),
		PriorityID:  uint(priorityID),
	}

	status, err := d.taskClient.CreateTask(session.Token, task)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=success&message=Tugas berhasil dibuat!")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
	}
}

func (d *dashboardWeb) TaskDeleteProcess(c *gin.Context) {
	// extracting email from context
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	// get session by email
	session, err := d.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
		return
	}

	// extracting task id from form value
	taskID, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message=Tugas tidak ditemukan!")
		return
	}

	// delete task
	status, err := d.taskClient.DeleteTask(session.Token, taskID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message="+err.Error())
	}

	if status == 200 {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=success&message=Berhasil menghapus tugas!")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message=Tugas gagal dihapus :(")
		return
	}
}