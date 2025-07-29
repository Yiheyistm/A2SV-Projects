package database

import (
	"errors"

	"github.com/yiheyistm/task_manager/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FromDomainToTaskEntity(u *domain.Task) (*TaskEntity, error) {
	if u == nil {
		return nil, errors.New("task cannot be nil")
	}

	return &TaskEntity{
		ID:          u.ID,
		Title:       u.Title,
		Description: u.Description,
		CreatedBy:   u.CreatedBy,
		DueDate:     primitive.NewDateTimeFromTime(u.DueDate),
		Status:      u.Status,
	}, nil
}

func FromTaskEntityToDomain(e *TaskEntity) *domain.Task {
	return &domain.Task{
		ID:          e.ID,
		CreatedBy:   e.CreatedBy,
		Title:       e.Title,
		Description: e.Description,
		DueDate:     e.DueDate.Time(),
		Status:      e.Status,
	}
}

func FromTaskEntityListToDomainList(entities []TaskEntity) []domain.Task {
	var tasks []domain.Task
	for _, entity := range entities {
		tasks = append(tasks, *FromTaskEntityToDomain(&entity))
	}
	return tasks
}

func FromStatusCountListToDomainList(entities []StatusCount) []domain.StatusCount {
	var statusCounts []domain.StatusCount
	for _, entity := range entities {
		statusCounts = append(statusCounts, domain.StatusCount{
			Status: entity.Status,
			Count:  entity.Count,
		})
	}
	return statusCounts
}
