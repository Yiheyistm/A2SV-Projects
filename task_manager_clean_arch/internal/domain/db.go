package domain

import "go.mongodb.org/mongo-driver/mongo"

var (
	DB       mongo.Database
	DBClient *mongo.Client
)
