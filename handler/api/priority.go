package api

import (
	"first-project/model"
	"first-project/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PriorityAPI interface {
	GetByID(c *gin.Context)
}

type priorityAPI struct {
	priorityService service.PriorityService
}

func NewPriorityAPI(priorityService service.PriorityService) *priorityAPI {
	return &priorityAPI{priorityService}
}

func (a *priorityAPI) GetByID(c *gin.Context) {
	priorityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid priority ID"))
		return
	}

	priority, err := a.priorityService.GetByID(priorityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, priority)
}