package web

import (
	"embed"
	"first-project/client"
	"first-project/service"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type AuthWeb interface {
	Register(c *gin.Context)
	RegisterProcess(c *gin.Context)
	Login(c *gin.Context)
	LoginProcess(c *gin.Context)
	Logout(c *gin.Context)
}

type authWeb struct {
	userClient     client.UserClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewAuthWeb(userClient client.UserClient, sessionService service.SessionService, embed embed.FS) *authWeb {
	return&authWeb{userClient, sessionService, embed}
}

func (a *authWeb) Register(c *gin.Context) {
	status := c.Query("status")
	message := c.Query("message")

	var dataTemplate = map[string]interface{}{
		"status":  status,
		"message": message,
	}

	var filepath = path.Join("views", "auth", "register.html")
	var header = path.Join("views", "general", "header.html")

	var temp, err = template.ParseFS(a.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/register?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/register?status=error&message="+err.Error())
		return
	}
}

func (a *authWeb) RegisterProcess(c *gin.Context) {
	fullname := c.Request.FormValue("fullname")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	status, err := a.userClient.Register(fullname, email, password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/register?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		c.Redirect(http.StatusSeeOther, "/client/login")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/register?status=error&message=Registrasi gagal!")
	}
}

func (a *authWeb) Login(c *gin.Context) {
	status := c.Query("status")
	message := c.Query("message")

	var dataTemplate = map[string]interface{}{
		"status":  status,
		"message": message,
	}

	var filepath = path.Join("views", "auth", "login.html")
	var header = path.Join("views", "general", "header.html")

	var temp, err = template.ParseFS(a.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/login?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/login?status=error&message="+err.Error())
		return
	}
}

func (a *authWeb) LoginProcess(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	status, err := a.userClient.Login(email, password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/login?status=error&message="+err.Error())
		return
	}

	session, err := a.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/login?status=error&message="+err.Error())
		return
	}

	if status == 200 {
		http.SetCookie(c.Writer, &http.Cookie{
			Name: "session_token",
			Value: session.Token,
			Path: "/",
			MaxAge: 31536000,
			Domain: "",
		})
		c.Redirect(http.StatusSeeOther, "/client/dashboard")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/login")
	}
}

func (a *authWeb) Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, false)
	c.Redirect(http.StatusSeeOther, "/client/login")
}