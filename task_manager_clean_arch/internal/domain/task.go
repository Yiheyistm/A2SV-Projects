package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date" validate:"required"`
	Status      string             `json:"status" bson:"status" validate:"oneof=pending completed"`
}

type StatusCount struct {
	Status string `bson:"_id"`
	Count  int    `bson:"count"`
}
