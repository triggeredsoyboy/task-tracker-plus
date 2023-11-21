package web

import (
	"embed"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type FeedbackWeb interface {
	FeedbackPage(c *gin.Context)
}

type feedbackWeb struct {
	embed embed.FS
}

func NewFeedbackWeb(embed embed.FS) *feedbackWeb {
	return &feedbackWeb{embed}
}

func (f *feedbackWeb) FeedbackPage(c *gin.Context) {
	var filepath = path.Join("views", "general", "feedback.html")
	var header = path.Join("viws", "general", "header.html")

	var temp = template.Must(template.ParseFS(f.embed, filepath, header))

	err := temp.Execute(c.Writer, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
}