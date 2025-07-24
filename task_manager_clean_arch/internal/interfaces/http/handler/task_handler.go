package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
)

type TaskHandler struct {
	TaskUsecase domain.ITaskUseCase
	UserUsecase domain.IUserUseCase
}

var validate = validator.New()

// List all the tasks
func (th *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := th.TaskUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, dto.FromDomainTaskToResponseList(tasks))
}

// Get a specific task by ID
func (th *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	task, err := th.TaskUsecase.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, dto.FromDomainTaskToResponse(&task))
}

// Create a specific task
func (th *TaskHandler) CreateTask(c *gin.Context) {
	user := th.UserUsecase.GetUserFromContext(c)
	var newTask dto.TaskRequest
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newTask.CreatedBy = user.Username
	err := th.TaskUsecase.Create(newTask.FromRequestToDomainTask())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, dto.FromDomainTaskToResponse(newTask.FromRequestToDomainTask()))
}

// Update a specific task by ID
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	updatedTask := &dto.TaskRequest{}
	if err := c.ShouldBindJSON(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := th.UserUsecase.GetUserFromContext(c)
	updatedTask.CreatedBy = user.Username
	err := th.TaskUsecase.Update(id, updatedTask.FromRequestToDomainTask())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, dto.FromDomainTaskToResponse(updatedTask.FromRequestToDomainTask()))
}

// Delete a specific task
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	err := th.TaskUsecase.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted"})

}

// Get task count by status
func (th *TaskHandler) GetTaskCountByStatus(c *gin.Context) {
	counts, err := th.TaskUsecase.GetTaskCountByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task counts"})
		return
	}
	c.JSON(http.StatusOK, counts)
}
