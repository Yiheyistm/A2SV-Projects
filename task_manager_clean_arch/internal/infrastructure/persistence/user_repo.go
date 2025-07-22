package persistence

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (s *userRepository) Insert(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user cannot be Empty")
	}
	_, err := s.database.Collection(s.collection).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	cursor, err := s.database.Collection(s.collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user domain.User
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

func (s *userRepository) getUser(ctx context.Context, key, value string) (*domain.User, error) {

	var user domain.User
	filter := bson.M{key: value}
	err := s.database.Collection(s.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.getUser(ctx, "username", username)

}

func (s *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.getUser(ctx, "email", email)

}

// func (s *userRepository) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
// 	jwtSecret := config.GetEnvString("JWT_SECRET", "my_secret_key")
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": user.Username,
// 		"role":     user.Role,
// 		"exp":      time.Now().Add(time.Hour * 48).Unix(),
// 	})
// 	tokenString, err := token.SignedString([]byte(jwtSecret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }

func (s *userRepository) GetUserFromContext(c *gin.Context) *domain.User {
	username := c.GetString("username")
	user, _ := s.GetByUsername(c, username)
	return user
}
