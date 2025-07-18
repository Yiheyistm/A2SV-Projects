package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	State   string `json:"state" bson:"state"`
	City    string `json:"city" bson:"city"`
	PinCode string `json:"pin_code" bson:"pin_code"`
}

type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"user_name"`
	Email   string             `json:"email" bson:"email"`
	Address Address
}
