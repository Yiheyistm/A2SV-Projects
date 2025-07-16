package controllers

import (
	"net/http"
	"strconv"
	"task_manager/models"
	"task_manager/services"

	"github.com/gin-gonic/gin"
)

// Safe handeler for / path
func SafeHandler(c *gin.Context) {
	tasks := services.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

// List all the tasks
func GetTasks(c *gin.Context) {
	tasks := services.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

// Get a specific task by ID
func GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	task, err := services.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Create a specific task
func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := services.Create(newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

// Update a specific task by ID
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existedTask, err := services.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	updatedTask := models.Task{}
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedTask.Id = existedTask.Id
	err = services.Update(updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// Delete a specific task
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existedTask, err := services.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	deleteTask := models.Task{}
	deleteTask.Id = existedTask.Id
	err = services.Delete(deleteTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted"})

}
