package main

import (
	"embed"
	"first-project/client"
	"first-project/db"
	"first-project/handler/api"
	"first-project/handler/web"
	"first-project/middleware"
	"first-project/model"
	repo "first-project/repository"
	"first-project/service"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	CategoryAPIHandler api.CategoryAPI
	PriorityAPIHandler api.PriorityAPI
	TaskAPIHandler     api.TaskAPI
}

type ClientHandler struct {
	FeedbackWeb  web.FeedbackWeb
	HomeWeb      web.HomeWeb
	AuthWeb      web.AuthWeb
	DashboardWeb web.DashboardWeb
	TaskWeb      web.TaskWeb
	CategoryWeb  web.CategoryWeb
}

//go:embed views/*
var Resources embed.FS

func main() {
	gin.SetMode(gin.DebugMode) // release mode

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()
		db := db.NewDB()
		router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s %s %s\"\n",
				params.TimeStamp.Format(time.RFC822),
				params.Method,
				params.Path,
				params.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())

		dbCredential := model.DBCredential{
			Host:         "localhost",
			Username:     "postgres",
			Password:     "dimaswicakk",
			DatabaseName: "test-web",
			Port:         5432,
			Schema:       "public",
		}

		connect, err := db.Connect(&dbCredential)
		if err != nil {
			panic("Gagal terhubung ke database")
		}

		connect.AutoMigrate(&model.User{}, &model.Session{}, &model.Category{}, &model.Priority{}, &model.Task{})

		router = RunServer(connect, router)
		router = RunClient(connect, router, Resources)

		fmt.Println("Server is running on port 8080")
		err = router.Run(":8080")
		if err != nil {
			panic("Router tidak dapat berjalan!")
		}
	}()

	wg.Wait()
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	sessionRepo := repo.NewSessionRepo(db)
	userRepo := repo.NewUserRepo(db)
	taskRepo := repo.NewTaskRepo(db)
	categoryRepo := repo.NewCategoryRepo(db)
	priorityRepo := repo.NewPriorityRepo(db)

	sessionService := service.NewSessionService(sessionRepo)
	userService := service.NewUserService(userRepo, sessionRepo, sessionService)
	taskService := service.NewTaskService(taskRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	priorityService := service.NewPriorityService(priorityRepo)

	userAPIHandler := api.NewUserAPI(userService)
	taskAPIHandler := api.NewTaskAPI(userService, categoryService, priorityService, taskService)
	categoryAPIHandler := api.NewCategoryAPI(userService, categoryService)
	priorityAPIHandler := api.NewPriorityAPI(priorityService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
		PriorityAPIHandler: priorityAPIHandler,
	}

	version := gin.Group("/api")
	{
		user := version.Group("/user")
		{
			user.POST("/register", apiHandler.UserAPIHandler.Register)
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.Use(middleware.Auth())
			user.GET("/profile", apiHandler.UserAPIHandler.GetCurrentUser)
		}

		task := version.Group("/task")
		{
			task.Use(middleware.Auth())
			task.POST("/add", apiHandler.TaskAPIHandler.CreateTask)
			task.GET("/list", apiHandler.TaskAPIHandler.TaskList)
			task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
			task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
			task.GET("/get/:id", apiHandler.TaskAPIHandler.GetByID)
			task.GET("/list-by-category", apiHandler.TaskAPIHandler.GetByCategory)
		}

		category := version.Group("/category")
		{
			category.Use(middleware.Auth())
			category.POST("/add", apiHandler.CategoryAPIHandler.CreateCategory)
			category.GET("/list", apiHandler.CategoryAPIHandler.CategoryList)
			category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
			category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
			category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetByID)

		}

		priority := version.Group("/priority")
		{
			priority.Use(middleware.Auth())
			priority.GET("/list", apiHandler.PriorityAPIHandler.GetByID)
		}

	}

	return gin
}

func RunClient(db *gorm.DB, gin *gin.Engine, embed embed.FS) *gin.Engine {
	sessionRepo := repo.NewSessionRepo(db)
	sessionService := service.NewSessionService(sessionRepo)

	userClient := client.NewUserClient()
	categoryClient := client.NewCategoryClient()
	taskClient := client.NewTaskClient()

	feedbackWeb := web.NewFeedbackWeb(embed)
	homeWeb := web.NewHomeWeb(embed)
	authWeb := web.NewAuthWeb(userClient, sessionService, embed)
	dashboardWeb := web.NewDashboardWeb(userClient, categoryClient, taskClient, sessionService, embed)
	taskWeb := web.NewTaskWeb(userClient, categoryClient, taskClient, sessionService, embed)
	categoryWeb := web.NewCategoryWeb(userClient, categoryClient, taskClient, sessionService, embed)

	client := ClientHandler{
		feedbackWeb, homeWeb, authWeb, dashboardWeb, taskWeb, categoryWeb,
	}

	// load css
	gin.StaticFS("/assets", http.Dir("assets"))
	gin.StaticFS("/node_modules", http.Dir("node_modules"))

	gin.GET("/", client.HomeWeb.Index)

	user := gin.Group("/client")
	{
		user.GET("/login", client.AuthWeb.Login)
		user.POST("/login/process", client.AuthWeb.LoginProcess)
		user.GET("/register", client.AuthWeb.Register)
		user.POST("/register/process", client.AuthWeb.RegisterProcess)
		user.Use(middleware.Auth())
		user.GET("/logout", client.AuthWeb.Logout)
	}

	main := gin.Group("/client")
	{
		main.Use(middleware.Auth())
		main.GET("/dashboard", client.DashboardWeb.Dashboard)
		main.POST("/dashboard/task-add/process", client.DashboardWeb.TaskAddProcess)
		main.POST("/dashboard/task-delete/process", client.DashboardWeb.TaskDeleteProcess)

		main.GET("/task/", client.TaskWeb.TaskPage)
		main.GET("/task/:id", client.TaskWeb.TaskByID)
		main.POST("/task/task-update/:id/process", client.TaskWeb.TaskUpdateProcess)

		main.GET("/category", client.CategoryWeb.CategoryPage)
		main.POST("/category/category-add/process", client.CategoryWeb.CategoryAddProcess)
		main.POST("/category/category-update/:id/process", client.CategoryWeb.CategoryUpdateProcess)
	}

	return gin
}