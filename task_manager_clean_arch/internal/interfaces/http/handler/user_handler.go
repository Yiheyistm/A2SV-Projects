package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserUsecase domain.UserUseCase
	TaskUsecase domain.TaskUseCase
	JwtService  *security.JwtService
}

func (uh *UserHandler) RegisterRequest(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	existedUser, _ := uh.UserUsecase.GetByUsername(strings.ToLower(newUser.Username))
	if existedUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := domain.User{
		Username: strings.ToLower(newUser.Username),
		Email:    strings.ToLower(newUser.Email),
		Password: string(hashPassword),
		Role:     newUser.Role,
	}

	err = uh.UserUsecase.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) LoginRequest(c *gin.Context) {
	var loginRequest domain.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if loginRequest.Identifier == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username/email and password are required"})
		return
	}
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, loginRequest.Identifier)

	var user *domain.User
	var err error

	if isEmail {
		user, err = uh.UserUsecase.GetByEmail(loginRequest.Identifier)
	} else {
		user, err = uh.UserUsecase.GetByUsername(loginRequest.Identifier)
	}
	fmt.Println("User:", user, "Error:", err)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	match := security.ValidatePassword(user.Password, loginRequest.Password)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	response, err := uh.JwtService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, domain.LoginResponse(response))
}

// TODO: Implement get all users, update user, get all user tasks, etc.

func (uh *UserHandler) GetAllUsers(c *gin.Context) {

	users, err := uh.UserUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	userName := c.Param("username")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is required"})
		return
	}

	user, err := uh.UserUsecase.GetByUsername(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserTasks
func (uh *UserHandler) GetUserTasks(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	fmt.Println("User from context:", user, username, c.GetString("role"))
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to see details about this user"})
		return
	}
	tasks, err := uh.TaskUsecase.GetTasksByUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks for user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetUserTask
func (uh *UserHandler) GetUserTask(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to see details about this user"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	task, err := uh.TaskUsecase.GetByIdAndUser(taskID, user.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateUserTask
func (uh *UserHandler) CreateUserTask(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create tasks on behalf of other user"})
		return
	}

	var newTask domain.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newTask.CreatedBy = user.Username

	err := uh.TaskUsecase.Create(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

// UpdateUserTask
func (uh *UserHandler) UpdateUserTask(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this task"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	updatedTask := &domain.Task{}
	if err := c.ShouldBindJSON(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedTask.CreatedBy = user.Username
	err := uh.TaskUsecase.UpdateByIdAndUser(taskID, updatedTask, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// DeleteUserTask
func (uh *UserHandler) DeleteUserTask(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete tasks on behalf of other user"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	err := uh.TaskUsecase.DeleteByIdAndUser(taskID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted successfully"})
}

// GetUserTaskStats
func (uh *UserHandler) GetUserTaskStats(c *gin.Context) {
	user := uh.UserUsecase.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to see details about this user"})
		return
	}

	stats, err := uh.TaskUsecase.GetTaskStatsByUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
