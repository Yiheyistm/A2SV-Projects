package routes

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/tasks", controllers.GetTasks)
		api.GET("/tasks/:id", controllers.GetTask)
		api.POST("/tasks", controllers.CreateTask)
		api.PUT("/tasks/:id", controllers.UpdateTask)
		api.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	return r
}
