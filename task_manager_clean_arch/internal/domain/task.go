package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID
	Title       string
	CreatedBy   string
	Description string
	DueDate     time.Time
	Status      string
}

type StatusCount struct {
	Status string `bson:"_id"`
	Count  int    `bson:"count"`
}

type TaskRepository interface {
	GetAll(context.Context) ([]Task, error)
	GetById(context.Context, string) (Task, error)
	GetByIdAndUser(context.Context, string, string) (Task, error)
	Create(context.Context, *Task) error
	Update(context.Context, string, *Task) error
	UpdateByIdAndUser(context.Context, string, *Task, string) error
	Delete(context.Context, string) error
	DeleteByIdAndUser(context.Context, string, string) error
	GetByUser(context.Context, string) ([]Task, error)
	GetTaskStatsByUser(context.Context, string) ([]StatusCount, error)
	GetTaskCountByStatus(context.Context) ([]StatusCount, error)
}

type ITaskUseCase interface {
	GetAll() ([]Task, error)
	GetById(string) (Task, error)
	GetByIdAndUser(string, string) (Task, error)
	Create(*Task) error
	Update(string, *Task) error
	UpdateByIdAndUser(string, *Task, string) error
	Delete(string) error
	DeleteByIdAndUser(string, string) error
	GetTasksByUser(string) ([]Task, error)
	GetTaskStatsByUser(string) ([]StatusCount, error)
	GetTaskCountByStatus() ([]StatusCount, error)
}
