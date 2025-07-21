package routes

import (
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/users/register", controllers.RegisterRequest)
		api.POST("/users/login", controllers.LoginRequest)
	}
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/users/:username/tasks", controllers.GetUserTasks)
		authGroup.GET("/users/:username/tasks/:id", controllers.GetUserTask)
		authGroup.POST("/users/:username/tasks", controllers.CreateUserTask)
		authGroup.PUT("/users/:username/tasks/:id", controllers.UpdateUserTask)
		authGroup.DELETE("/users/:username/tasks/:id", controllers.DeleteUserTask)
		authGroup.GET("/users/:username/tasks/stats", controllers.GetUserTaskStats)
	}
	adminGroup := authGroup.Group("/")
	adminGroup.Use(middleware.AdminOnlyMiddleware())
	{
		adminGroup.GET("/tasks", controllers.GetTasks)
		adminGroup.GET("/tasks/stats", controllers.GetTaskCountByStatus)
		adminGroup.GET("/tasks/:id", controllers.GetTask)
		adminGroup.POST("/tasks", controllers.CreateTask)
		adminGroup.PUT("/tasks/:id", controllers.UpdateTask)
		adminGroup.DELETE("/tasks/:id", controllers.DeleteTask)

		adminGroup.GET("/users", controllers.GetAllUsers)
		adminGroup.GET("/users/:username", controllers.GetUser)
	}

	return r
}
