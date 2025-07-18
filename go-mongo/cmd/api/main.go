package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yiheyistm/go-mongo/config"
	"github.com/yiheyistm/go-mongo/controllers"
	"github.com/yiheyistm/go-mongo/routes"
	"github.com/yiheyistm/go-mongo/services"
)

func main() {
	port := config.GetEnvString("PORT", "8080")
	url := config.GetEnvString("MONGO_URI", "mongodb://localhost:27017")
	db := config.GetEnvString("MONGO_DB_NAME", "")
	collName := config.GetEnvString("MONGO_COLLECTION_NAME", "users")
	collection := config.ConnectDB(url, db, collName)
	userService := services.NewUserService(collection, context.TODO())
	userController := controllers.NewUserController(userService)
	r := routes.SetupRoutes(&userController)

	r.Run(":" + port) // listen and serve on
}
