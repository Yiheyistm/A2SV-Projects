package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password" bson:"password" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" bson:"identifier" validate:"required"`
	Password   string `json:"password" bson:"password" validate:"required,min=6"`
}
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

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
