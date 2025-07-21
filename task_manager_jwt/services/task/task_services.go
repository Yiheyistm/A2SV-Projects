package services

import (
	"task_manager/models"
)

type TaskServices interface {
	GetAll() ([]models.Task, error)
	GetById(string) (models.Task, error)
	GetByIdAndUser(string, string) (models.Task, error)
	Create(*models.Task) error
	Update(string, *models.Task) error
	UpdateByIdAndUser(string, *models.Task, string) error
	Delete(string) error
	DeleteByIdAndUser(string, string) error
	GetByUser(string) ([]models.Task, error)
	GetTaskStatsByUser(string) ([]models.StatusCount, error)
	GetTaskCountByStatus() ([]models.StatusCount, error)
}
