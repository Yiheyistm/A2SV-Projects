package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/go-mongo/controllers"
)

func SetupRoutes(userController *controllers.UserController) *gin.Engine {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	{
		v1.GET("/users", userController.GetAllUsers)
		v1.GET("/users/:name", userController.GetUser)
		v1.POST("/users", userController.CreateUser)
		v1.PATCH("/users/:name", userController.UpdateUser)
		v1.DELETE("/users/:name", userController.DeleteUser)

	}
	return g
}
