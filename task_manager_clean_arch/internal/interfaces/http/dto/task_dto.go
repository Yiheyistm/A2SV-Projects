package dto

import (
	"time"
)

type TaskRequest struct {
	Title       string    `json:"title" bson:"title" validate:"required"`
	CreatedBy   string    `json:"created_by" bson:"created_by"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date" validate:"required"`
	Status      string    `json:"status" bson:"status" validate:"oneof=pending completed"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	CreatedBy   string    `json:"created_by"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}
