package services

import (
	"errors"
	"task_manager/models"
	"time"
)

var tasks = []models.Task{
	{Id: 1, Title: "Task One", Description: "Description for Task One", DueDate: time.Now().Add(24 * time.Hour), Status: "pending"},
	{Id: 2, Title: "Task Two", Description: "Description for Task Two", DueDate: time.Now().Add(48 * time.Hour), Status: "completed"},
	{Id: 3, Title: "Task Three", Description: "Description for Task Three", DueDate: time.Now().Add(72 * time.Hour), Status: "pending"},
}

func GetAllTasks() []models.Task {
	return tasks
}

func GetById(id int) (models.Task, error) {

	for _, task := range tasks {
		if task.Id == id {
			return task, nil
		}
	}
	return models.Task{}, errors.New("task not found")
}

func Create(task models.Task) error {
	task.Id = len(tasks) + 1
	tasks = append(tasks, task)
	return nil
}

func Update(task models.Task) error {
	tasks[task.Id] = task
	return nil
}

func Delete(task models.Task) error {
	tasks = append(tasks[:task.Id], tasks[task.Id+1:]...)
	return nil
}
