package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/interfaces/middleware"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(env *config.Env, db mongo.Database) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware(env.AccessTokenSecret))
	adminGroup := authGroup.Group("/")
	adminGroup.Use(middleware.AdminOnlyMiddleware())

	AuthRoutes(env, db, api)
	UserRoutes(env, db, authGroup)
	TaskRoutes(env, db, adminGroup)

	return r
}
