package models

import (
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
	Token string `json:"token" bson:"token"`
}
