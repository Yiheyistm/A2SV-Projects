package services

import (
	"context"
	"log"

	"github.com/yiheyistm/go-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServicesImpl struct {
	UserCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) UserServices {
	return &UserServicesImpl{
		UserCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServicesImpl) GetAll() ([]models.User, error) {
	var users []models.User
	cursor, err := u.UserCollection.Find(u.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(u.ctx)

	for cursor.Next(u.ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserServicesImpl) Create(user *models.User) error {
	_, err := u.UserCollection.InsertOne(u.ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServicesImpl) GetSingle(name *string) (models.User, error) {
	var user models.User
	filter := bson.D{{Key: "user_name", Value: *name}}
	err := u.UserCollection.FindOne(u.ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}
func (u *UserServicesImpl) Update(user *models.User) error {
	name := user.Name
	filter := bson.M{"user_name": name}
	update := bson.M{"$set": bson.M{
		"user_name": user.Name,
		"email":     user.Email,
		"address":   user.Address,
	}}
	result, err := u.UserCollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update user: %v\n", err)
		return err
	}
	log.Printf("Updated %v documents\n", result.ModifiedCount)
	return nil
}

func (u *UserServicesImpl) Delete(name *string) error {
	result, err := u.UserCollection.DeleteOne(u.ctx, bson.M{"user_name": *name})
	if err != nil {
		log.Printf("Failed to delete user: %v\n", err)
		return err
	}
	log.Printf("Deleted %v documents\n", result.DeletedCount)
	return nil
}
