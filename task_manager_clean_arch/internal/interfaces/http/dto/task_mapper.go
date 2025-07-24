package dto

import (
	"github.com/yiheyistm/task_manager/internal/domain"
)

func (r *TaskRequest) FromRequestToDomainTask() *domain.Task {
	return &domain.Task{
		Title:       r.Title,
		CreatedBy:   r.CreatedBy,
		Description: r.Description,
		DueDate:     r.DueDate,
		Status:      r.Status,
	}
}
func FromDomainTaskToResponse(task *domain.Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		CreatedBy:   task.CreatedBy,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}

func FromDomainTaskToResponseList(tasks []domain.Task) []TaskResponse {
	var taskResponses []TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, *FromDomainTaskToResponse(&task))
	}
	return taskResponses
}
