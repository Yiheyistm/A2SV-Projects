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

func AuthRoutes(env *config.Env, db mongo.Database, group *gin.RouterGroup) {
	ur := persistence.NewUserRepository(db, env.DBUserCollection)
	tr := persistence.NewTaskRepository(db, env.DBTaskCollection)
	jwtService := security.NewJWTService(
		env.AccessTokenSecret,
		env.RefreshTokenSecret,
		env.AccessTokenExpiryHour,
		env.RefreshTokenExpiryHour,
	)
	userHandler := handler.UserHandler{
		JwtService:  jwtService,
		TaskUsecase: usecase.NewTaskUseCase(tr),
		UserUsecase: usecase.NewUserUseCase(ur),
	}
	group.POST("/users/register", userHandler.RegisterRequest)
	group.POST("/users/login", userHandler.LoginRequest)
}
