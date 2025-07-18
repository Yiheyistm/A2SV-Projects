package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/go-mongo/models"
	"github.com/yiheyistm/go-mongo/services"
)

type UserController struct {
	UserServices services.UserServices
}

func NewUserController(userServices services.UserServices) UserController {
	return UserController{
		UserServices: userServices,
	}
}
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.UserServices.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (uc *UserController) GetUser(c *gin.Context) {
	userName := c.Param("name")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is required"})
		return
	}

	user, err := uc.UserServices.GetSingle(&userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if (user == models.User{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})

}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := uc.UserServices.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userName := c.Param("name")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is required"})
		return
	}

	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user.Name = userName
	if err := uc.UserServices.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userName := c.Param("name")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name is required"})
		return
	}

	if err := uc.UserServices.Delete(&userName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
