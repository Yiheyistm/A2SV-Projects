package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yiheyistm/task_manager/internal/domain"
)

type JwtService struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessExpiry, refreshExpiry int) *JwtService {
	return &JwtService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  time.Duration(accessExpiry) * time.Hour,
		refreshExpiry: time.Duration(refreshExpiry) * time.Hour,
	}
}

func (s *JwtService) GenerateTokens(user domain.User) (domain.RefreshToken, error) {
	fmt.Println(s, user)
	accessClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(s.accessExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return domain.RefreshToken{}, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.refreshExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return domain.RefreshToken{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *JwtService) ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *JwtService) ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.refreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
