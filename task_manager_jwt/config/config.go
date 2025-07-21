package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(uri, dbName, collName string) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	// Ping to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	collection := client.Database(dbName).Collection(collName)

	return collection
}

func CreateTaskIndexes(coll *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}},
		Options: options.Index().SetName("status_index"),
	}

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Index creation failed:", err)
	} else {
		log.Println("Index created on status field")
	}
}

func CreateUserIndexes(coll *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_index"),
	}

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Index creation failed:", err)
	} else {
		log.Println("Index created on username field")
	}
}
