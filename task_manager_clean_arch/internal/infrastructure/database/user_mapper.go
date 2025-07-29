package database

import (
	"errors"

	"github.com/yiheyistm/task_manager/internal/domain"
)

func FromDomainToEntity(u *domain.User) (*UserEntity, error) {
	if u == nil {
		return nil, errors.New("user cannot be nil")
	}
	return &UserEntity{
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
