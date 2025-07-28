package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yiheyistm/task_manager/internal/domain"
)

type refreshTokenUsecase struct {
	userRepository   domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
}

func NewRefreshTokenUsecase(userRepository domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository) domain.IRefreshTokenUsecase {
	return &refreshTokenUsecase{
		userRepository:   userRepository,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (rtu *refreshTokenUsecase) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rtu.userRepository.GetByUsername(ctx, username)
}

func (rtu *refreshTokenUsecase) GenerateTokens(user domain.User) (domain.RefreshToken, error) {
	return rtu.refreshTokenRepo.GenerateTokens(user)
}

func (rtu *refreshTokenUsecase) ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	return rtu.refreshTokenRepo.ValidateRefreshToken(token)
}
func (rtu *refreshTokenUsecase) ValidateToken(token string) (jwt.MapClaims, error) {
	return rtu.refreshTokenRepo.ValidateToken(token)
}
