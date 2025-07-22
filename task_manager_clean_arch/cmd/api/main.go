package main

import (
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/infrastructure/database"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/router"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	app := App()
	env := app.Env
	db := *app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()
	route := router.SetupRouter(env, db)
	route.Run(env.ServerAddress)
}

type Application struct {
	Env   *config.Env
	Mongo mongo.Client
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
