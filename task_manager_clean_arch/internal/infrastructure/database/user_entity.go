package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserEntity struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password,omitempty"`
	Role     string             `bson:"role"`
}
