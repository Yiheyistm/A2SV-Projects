package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenUsecase interface {
	GenerateTokens(user User) (RefreshTokenResponse, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(token string) (jwt.MapClaims, error)
	GetByUsername(string) (*User, error)
}
