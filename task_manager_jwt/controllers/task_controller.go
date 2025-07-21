package controllers

import (
	"net/http"
	"task_manager/models"
	services "task_manager/services/task"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var taskServices = services.NewTaskServicesImpl()
var validate = validator.New()

// List all the tasks
func GetTasks(c *gin.Context) {
	tasks, err := taskServices.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Get a specific task by ID
func GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	task, err := taskServices.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Create a specific task
func CreateTask(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newTask.CreatedBy = user.Username
	err := taskServices.Create(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

// Update a specific task by ID
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	updatedTask := &models.Task{}
	if err := c.ShouldBindJSON(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := userServices.GetUserFromContext(c)
	updatedTask.CreatedBy = user.Username
	err := taskServices.Update(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// Delete a specific task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	err := taskServices.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted"})

}

// Get task count by status
func GetTaskCountByStatus(c *gin.Context) {
	counts, err := taskServices.GetTaskCountByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task counts"})
		return
	}
	c.JSON(http.StatusOK, counts)
}
