package domain

import (
	"context"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string
	Username string
	Email    string
	Password string
	Role     string
}
type UserRepository interface {
	GetAll(context.Context) ([]User, error)
	GetByUsername(context.Context, string) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	Insert(context.Context, *User) error
	GetUser(context.Context, string, string) (*User, error)
	GetUserFromContext(c *gin.Context) *User
}

type IUserUseCase interface {
	GetAll() ([]User, error)
	GetByUsername(string) (*User, error)
	GetByEmail(string) (*User, error)
	Insert(*User) error
	// GenerateToken(*User) (string, error)
	GetUserFromContext(*gin.Context) *User
}
