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

type CategoryWeb interface {
	CategoryPage(c *gin.Context)
	CategoryAddProcess(c *gin.Context)
	CategoryUpdateProcess(c *gin.Context)
}

type categoryWeb struct {
	userClient     client.UserClient
	categoryClient client.CategoryClient
	taskClient     client.TaskClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewCategoryWeb(userClient client.UserClient, categoryClient client.CategoryClient, taskClient client.TaskClient, sessionService service.SessionService, embed embed.FS) *categoryWeb {
	return &categoryWeb{userClient, categoryClient, taskClient, sessionService, embed}
}

func (cat *categoryWeb) CategoryPage(c *gin.Context) {
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
	session, err := cat.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	// get the current user's details
	userDetail, err := cat.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	// get the current user's categories
	userCategories, err := cat.categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	userTasks, err := cat.taskClient.TaskByCategory(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	var dataTemplate = map[string]interface{}{
		"status":          status,
	    "message":         message,
		"user_detail":     userDetail,
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

	var filepath = path.Join("views", "main", "category.html")
	var header = path.Join("views", "general", "header.html")

	temp, err := template.New("category.html").Funcs(funcMap).ParseFS(cat.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
	}
}

func (cat *categoryWeb) CategoryAddProcess(c *gin.Context) {
	// extracting email from context
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	// get session by email
	session, err := cat.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	// get the current user's details
	userDetail, err := cat.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	category := model.Category{
		Name:   c.Request.FormValue("name"),
		UserID: userDetail.ID,
	}

	status, err := cat.categoryClient.CreateCategory(session.Token, category)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		c.Redirect(http.StatusSeeOther, "/client/category?status=success&message=Berhasil membuat kategori!")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message+"+err.Error())
	}
}

func (cat *categoryWeb) CategoryUpdateProcess(c *gin.Context) {
	// extracting email from context
	var email string
	if data, ok := c.Get("email"); ok {
		if contextData, ok := data.(string); ok {
			email = contextData
		}
	}

	// get session by email
	session, err := cat.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	// get the current user's details
	userDetail, err := cat.userClient.GetCurrentUser(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	// extracting category id from query parameters
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message=Category ID is missing")
		return
	}

	// convert category id to integer
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message"+err.Error())
		return
	}

	category := model.Category{
		Name:   c.Request.FormValue("name"),
		UserID: userDetail.ID,
	}

	status, err := cat.categoryClient.UpdateCategory(session.Token, categoryID, category)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/category?status=error&message="+err.Error())
		return
	}

	if status == 200 {
		c.Redirect(http.StatusSeeOther, "/client/category?status=success&message=Berhasil mengubah kategori!")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/dashboard?status=error&message+"+err.Error())
	}
}