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

func UserRoutes(env *config.Env, db mongo.Database, protectedGroup *gin.RouterGroup, adminGroup *gin.RouterGroup) {
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
	adminGroup.GET("/users", userHandler.GetAllUsers)
	adminGroup.GET("/users/:username", userHandler.GetUser)
	protectedGroup.GET("/users/:username/tasks", userHandler.GetUserTasks)
	protectedGroup.GET("/users/:username/tasks/:id", userHandler.GetUserTask)
	protectedGroup.POST("/users/:username/tasks", userHandler.CreateUserTask)
	protectedGroup.PUT("/users/:username/tasks/:id", userHandler.UpdateUserTask)
	protectedGroup.DELETE("/users/:username/tasks/:id", userHandler.DeleteUserTask)
	protectedGroup.GET("/users/:username/tasks/stats", userHandler.GetUserTaskStats)
}
