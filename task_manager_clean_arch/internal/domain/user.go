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

// type LoginRequest struct {
// 	Identifier string
// 	Password   string
// }
// type LoginResponse struct {
// 	AccessToken  string
// 	RefreshToken string
// }

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Insert(ctx context.Context, user *User) error
	// GenerateToken(ctx context.Context, user *User) (string, error)
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
