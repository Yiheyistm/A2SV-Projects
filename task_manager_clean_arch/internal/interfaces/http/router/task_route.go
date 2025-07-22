package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/infrastructure/persistence"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/handler"
	"github.com/yiheyistm/task_manager/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func TaskRoutes(env *config.Env, db mongo.Database, group *gin.RouterGroup) {
	tr := persistence.NewTaskRepository(db, env.DBTaskCollection)
	ur := persistence.NewUserRepository(db, env.DBUserCollection)
	taskHandler := handler.TaskHandler{
		TaskUsecase: usecase.NewTaskUseCase(tr),
		UserUsecase: usecase.NewUserUseCase(ur),
	}
	group.GET("/tasks", taskHandler.GetTasks)
	group.GET("/tasks/stats", taskHandler.GetTaskCountByStatus)
	group.GET("/tasks/:id", taskHandler.GetTask)
	group.POST("/tasks", taskHandler.CreateTask)
	group.PUT("/tasks/:id", taskHandler.UpdateTask)
	group.DELETE("/tasks/:id", taskHandler.DeleteTask)
}
