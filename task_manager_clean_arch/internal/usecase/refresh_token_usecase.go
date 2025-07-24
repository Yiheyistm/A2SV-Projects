package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
)

type refreshTokenUsecase struct {
	userRepository domain.UserRepository
	jwtService     security.JwtService
}

func NewRefreshTokenUsecase(userRepository domain.UserRepository, jwtService security.JwtService) domain.IRefreshTokenUsecase {
	return &refreshTokenUsecase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (rtu *refreshTokenUsecase) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rtu.userRepository.GetByUsername(ctx, username)
}

func (rtu *refreshTokenUsecase) GenerateTokens(user domain.User) (domain.RefreshToken, error) {
	return rtu.jwtService.GenerateTokens(user)
}

func (rtu *refreshTokenUsecase) ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	return rtu.jwtService.ValidateRefreshToken(token)
}
func (rtu *refreshTokenUsecase) ValidateToken(token string) (jwt.MapClaims, error) {
	return rtu.jwtService.ValidateToken(token)
}
