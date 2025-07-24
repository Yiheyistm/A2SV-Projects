package database

import (
	"github.com/yiheyistm/task_manager/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FromDomainToEntity(u *domain.User) (*UserEntity, error) {
	objectID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, err
	}
	return &UserEntity{
		ID:       objectID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}, nil
}

func FromEntityToDomain(e *UserEntity) *domain.User {
	return &domain.User{
		ID:       e.ID.Hex(),
		Username: e.Username,
		Email:    e.Email,
		Password: e.Password,
		Role:     e.Role,
	}
}
func FromEntityListToDomainList(entities []UserEntity) []domain.User {
	var users []domain.User
	for _, entity := range entities {
		users = append(users, *FromEntityToDomain(&entity))
	}
	return users
}
