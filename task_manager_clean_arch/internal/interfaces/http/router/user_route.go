package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/infrastructure/persistence"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/handler"
	"github.com/yiheyistm/task_manager/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(env *config.Env, db mongo.Database, group *gin.RouterGroup) {
	ur := persistence.NewUserRepository(db, env.DBUserCollection)
	tr := persistence.NewTaskRepository(db, env.DBTaskCollection)
	jwtService := security.NewJWTService(env.AccessTokenSecret, time.Duration(env.AccessTokenExpiryHour))
	userHandler := handler.UserHandler{
		JwtService:  jwtService,
		TaskUsecase: usecase.NewTaskUseCase(tr),
		UserUsecase: usecase.NewUserUseCase(ur),
	}
	group.GET("/users/:username/tasks", userHandler.GetUserTasks)
	group.GET("/users/:username/tasks/:id", userHandler.GetUserTask)
	group.POST("/users/:username/tasks", userHandler.CreateUserTask)
	group.PUT("/users/:username/tasks/:id", userHandler.UpdateUserTask)
	group.DELETE("/users/:username/tasks/:id", userHandler.DeleteUserTask)
	group.GET("/users/:username/tasks/stats", userHandler.GetUserTaskStats)
}
