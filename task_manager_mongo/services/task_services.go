package services

import (
	"task_manager/models"
)

type TaskServices interface {
	GetAll() ([]models.Task, error)
	GetById(id string) (models.Task, error)
	Create(task *models.Task) error
	Update(id string, updateTask *models.Task) error
	Delete(id string) error
}
