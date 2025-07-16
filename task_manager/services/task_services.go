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

func GetById(id int) (int, models.Task, error) {

	for index, task := range tasks {
		if task.Id == id {
			return index, task, nil
		}
	}
	return -1, models.Task{}, errors.New("task not found")
}

func Create(task *models.Task) error {
	task.Id = tasks[len(tasks)-1].Id + 1
	tasks = append(tasks, *task)
	return nil
}

func Update(index int, updateTask *models.Task) error {
	if index < 0 || index >= len(tasks) {
		return errors.New("task not found")
	}
	tasks[index] = *updateTask
	return nil
}

func Delete(index int) error {
	if index < 0 || index >= len(tasks) {
		return errors.New("task not found")
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	return nil
}
