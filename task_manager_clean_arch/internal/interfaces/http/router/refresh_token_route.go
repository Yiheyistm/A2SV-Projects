package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/infrastructure/persistence"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/handler"
	"github.com/yiheyistm/task_manager/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func RefreshTokenRoutes(env *config.Env, db mongo.Database, group *gin.RouterGroup) {
	ur := persistence.NewUserRepository(db, env.DBUserCollection)
	jwtService := security.NewJWTService(
		env.AccessTokenSecret,
		env.RefreshTokenSecret,
		env.AccessTokenExpiryHour,
		env.RefreshTokenExpiryHour,
	)
	userHandler := handler.RefreshTokenHandler{
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(ur, *jwtService),
	}
	group.POST("/users/refresh", userHandler.RefreshToken)
}
