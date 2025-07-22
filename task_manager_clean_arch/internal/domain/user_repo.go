package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	// "github.com/yiheyistm/task_manager/internal/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Insert(ctx context.Context, user *User) error
	// GenerateToken(ctx context.Context, user *User) (string, error)
	GetUserFromContext(c *gin.Context) *User
}

type UserUseCase interface {
	GetAll() ([]User, error)
	GetByUsername(string) (*User, error)
	GetByEmail(string) (*User, error)
	Insert(*User) error
	// GenerateToken(*User) (string, error)
	GetUserFromContext(*gin.Context) *User
}
