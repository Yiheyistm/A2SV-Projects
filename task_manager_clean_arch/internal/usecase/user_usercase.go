package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) GetAll() ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	users, err := uc.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *UserUseCase) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (uc *UserUseCase) GetByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (uc *UserUseCase) Insert(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if user == nil {
		return errors.New("user cannot be nil")
	}
	err := uc.userRepo.Insert(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) GetUserFromContext(c *gin.Context) *domain.User {
	user := uc.userRepo.GetUserFromContext(c)
	if user == nil {
		return &domain.User{}
	}
	return user
}
