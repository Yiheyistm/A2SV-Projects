package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yiheyistm/task_manager/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(env *config.Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DBHostURI := env.DBHostURI
	if DBHostURI == "" {
		DBHostURI = fmt.Sprintf("mongodb://%s:%s", env.DBHost, env.DBPort)
	}

	clientOptions := options.Client().ApplyURI(DBHostURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return *client
}

func CloseMongoDBConnection(client mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
