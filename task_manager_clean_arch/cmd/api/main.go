package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/infrastructure/database"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/router"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	app := App()
	env := app.Env
	if env.AppEnv == "development" {
		// make development logs
		gin.SetMode(gin.DebugMode)
		fmt.Println("------------------------ Development Mode ------------------------")
	} else {
		// make production logs
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("------------------------ Production Mode ------------------------")
	}
	db := *app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	route := router.SetupRouter(env, db)
	route.Run(env.ServerAddress)
}

type Application struct {
	Env   *config.Env
	Mongo *mongo.Client
}

func App() Application {
	app := &Application{}
	app.Env = config.Load()
	app.Mongo = database.NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	database.CloseMongoDBConnection(app.Mongo)
}
