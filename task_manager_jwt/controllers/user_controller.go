package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"task_manager/models"
	services "task_manager/services/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var userServices = services.NewUserServicesImpl()

func RegisterRequest(c *gin.Context) {

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	existedUser, _ := userServices.GetByUsername(strings.ToLower(newUser.Username))
	if existedUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{
		Username: strings.ToLower(newUser.Username),
		Email:    strings.ToLower(newUser.Email),
		Password: string(hashPassword),
		Role:     newUser.Role,
	}

	err = userServices.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func LoginRequest(c *gin.Context) {
	var loginRequest models.LoginRequest
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

	var user *models.User
	var err error

	if isEmail {
		user, err = userServices.GetByEmail(loginRequest.Identifier)
	} else {
		user, err = userServices.GetByUsername(loginRequest.Identifier)
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something went wrong"})
		return
	}

	token, err := userServices.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

// TODO: Implement get all users, update user, get all user tasks, etc.

func GetAllUsers(c *gin.Context) {

	users, err := userServices.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUser(c *gin.Context) {
	userName := c.Param("username")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is required"})
		return
	}

	user, err := userServices.GetByUsername(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserTasks
func GetUserTasks(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
	username := c.Param("username")
	fmt.Println("User from context:", user, username, c.GetString("role"))
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to see details about this user"})
		return
	}
	tasks, err := taskServices.GetByUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks for user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// Getusertask
func GetUserTask(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
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

	task, err := taskServices.GetByIdAndUser(taskID, user.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateUserTask
func CreateUserTask(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create tasks on behalf of other user"})
		return
	}

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

// UpdateUserTask
func UpdateUserTask(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
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

	updatedTask := &models.Task{}
	if err := c.ShouldBindJSON(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedTask.CreatedBy = user.Username
	err := taskServices.UpdateByIdAndUser(taskID, updatedTask, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// DeleteUserTask
func DeleteUserTask(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
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

	err := taskServices.DeleteByIdAndUser(taskID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted successfully"})
}

// GetUserTaskStats
func GetUserTaskStats(c *gin.Context) {
	user := userServices.GetUserFromContext(c)
	username := c.Param("username")
	if user.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to see details about this user"})
		return
	}

	stats, err := taskServices.GetTaskStatsByUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
