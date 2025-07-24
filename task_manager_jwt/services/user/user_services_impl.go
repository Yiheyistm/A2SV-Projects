package services

import (
	"context"
	"errors"
	"task_manager/config"
	"task_manager/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServicesImpl struct {
	UserCollection *mongo.Collection
}

func NewUserServicesImpl() UserServices {
	url := config.GetEnvString("MONGO_URI", "mongodb://localhost:27017")
	db := config.GetEnvString("MONGO_DB_NAME", "")
	collName := config.GetEnvString("MONGO_COLLECTION_USERS", "users")
	collection := config.ConnectDB(url, db, collName)
	config.CreateUserIndexes(collection)
	return &UserServicesImpl{
		UserCollection: collection,
	}
}

func (s *UserServicesImpl) Insert(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if user == nil {
		return errors.New("user cannot be Empty")
	}
	_, err := s.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServicesImpl) GetAll() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var users []models.User
	cursor, err := s.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
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

func (s *UserServicesImpl) getUser(key, value string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{key: value}
	err := s.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserServicesImpl) GetByUsername(username string) (*models.User, error) {
	return s.getUser("username", username)

}

func (s *UserServicesImpl) GetByEmail(email string) (*models.User, error) {
	return s.getUser("email", email)

}

func (s *UserServicesImpl) GenerateToken(user *models.User) (string, error) {
	jwtSecret := config.GetEnvString("JWT_SECRET", "my_secret_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserServicesImpl) GetUserFromContext(c *gin.Context) *models.User {
	username := c.GetString("username")
	user, _ := s.GetByUsername(username)
	return user
}
