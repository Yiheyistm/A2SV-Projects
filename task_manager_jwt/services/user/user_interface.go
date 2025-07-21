package services

import (
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

type UserServices interface {
	GetAll() ([]models.User, error)
	GetByUsername(string) (*models.User, error)
	GetByEmail(string) (*models.User, error)
	Insert(*models.User) error
	GenerateToken(*models.User) (string, error)
	GetUserFromContext(*gin.Context) *models.User
}
