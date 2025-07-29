// TAsk entities
package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	CreatedBy   string             `bson:"created_by"`
	Description string             `bson:"description"`
	DueDate     primitive.DateTime `bson:"due_date"`
	Status      string             `bson:"status"`
}

type StatusCount struct {
	Status string `bson:"_id"`
	Count  int    `bson:"count"`
}
