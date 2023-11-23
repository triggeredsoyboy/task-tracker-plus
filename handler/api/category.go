package api

import (
	"first-project/model"
	"first-project/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryAPI interface {
	CreateCategory(c *gin.Context)
	CategoryList(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	GetByID(c *gin.Context)
}

type categoryAPI struct {
	userService     service.UserService
	categoryService service.CategoryService
}

func NewCategoryAPI(userService service.UserService, categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{userService, categoryService}
}

func (a *categoryAPI) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryReq
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

	reqBody := model.Category{
		Name: req.Name,
		UserID: req.UserID,
	}

	err := a.categoryService.CreateCategory(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "category created",
		"data":    reqBody,
	})
}

func (a *categoryAPI) CategoryList(c *gin.Context) {
	currentUser := a.userService.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("user not authenticated"))
		return
	}

	categories, err := a.categoryService.CategoryList(int(currentUser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (a *categoryAPI) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	var updatedCategory model.UpdateCategoryReq
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	err = a.categoryService.UpdateCategory(categoryID, &updatedCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewErrorResponse("update category success"))
}

func (a *categoryAPI) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	err = a.categoryService.DeleteCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewErrorResponse("delete category success"))
}

func (a *categoryAPI) GetByID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	category, err := a.categoryService.GetByID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, category)
}