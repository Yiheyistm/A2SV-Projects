package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
)

// JwtServiceSuite defines the test suite for JwtService
type JwtServiceSuite struct {
	suite.Suite
	jwtService *security.JwtService
}

// SetupTest initializes the JwtService before each test
func (s *JwtServiceSuite) SetupTest() {
	s.jwtService = security.NewJWTService("access_secret", "refresh_secret", 1, 24).(*security.JwtService)
}

// TestJwtServiceSuite runs the test suite
func TestJwtServiceSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceSuite))
}

// TestNewJWTService tests the NewJWTService function
func (s *JwtServiceSuite) TestNewJWTService() {
	s.Run("Success", func() {
		accessSecret := "access_secret"
		refreshSecret := "refresh_secret"
		accessExpiry := 1   // 1 hour
		refreshExpiry := 24 // 24 hours

		service := security.NewJWTService(accessSecret, refreshSecret, accessExpiry, refreshExpiry).(*security.JwtService)

		s.Equal(accessSecret, service.AccessSecret)
		s.Equal(refreshSecret, service.RefreshSecret)
		s.Equal(time.Duration(accessExpiry)*time.Hour, service.AccessExpiry)
		s.Equal(time.Duration(refreshExpiry)*time.Hour, service.RefreshExpiry)
	})

	s.Run("EmptySecrets", func() {
		service := security.NewJWTService("", "", 1, 24).(*security.JwtService)

		s.Empty(service.AccessSecret)
		s.Empty(service.RefreshSecret)
		s.Equal(time.Duration(1)*time.Hour, service.AccessExpiry)
		s.Equal(time.Duration(24)*time.Hour, service.RefreshExpiry)
	})
}

// TestGenerateTokens tests the GenerateTokens function
func (s *JwtServiceSuite) TestGenerateTokens() {
	s.Run("Success", func() {
		user := domain.User{
			ID:       "1",
			Username: "Abebe",
			Role:     "user",
		}

		tokens, err := s.jwtService.GenerateTokens(user)

		s.NoError(err)
		s.NotEmpty(tokens.AccessToken)
		s.NotEmpty(tokens.RefreshToken)

		// Validate access token
		accessClaims, err := s.jwtService.ValidateToken(tokens.AccessToken)
		s.NoError(err)
		s.Equal(user.ID, accessClaims["sub"])
		s.Equal(user.Username, accessClaims["username"])
		s.Equal(user.Role, accessClaims["role"])
		s.NotEmpty(accessClaims["iat"])
		s.NotEmpty(accessClaims["exp"])

		// Validate refresh token
		refreshClaims, err := s.jwtService.ValidateRefreshToken(tokens.RefreshToken)
		s.NoError(err)
		s.Equal(user.ID, refreshClaims["sub"])
		s.Equal(user.Username, refreshClaims["username"])
		s.NotEmpty(refreshClaims["iat"])
		s.NotEmpty(refreshClaims["exp"])
	})

	s.Run("EmptyUserID", func() {
		user := domain.User{
			ID:       "",
			Username: "Abebe",
			Role:     "user",
		}

		tokens, err := s.jwtService.GenerateTokens(user)

		s.NoError(err) // JWT allows empty sub
		s.NotEmpty(tokens.AccessToken)
		s.NotEmpty(tokens.RefreshToken)
	})
}

// TestValidateToken tests the ValidateToken function
func (s *JwtServiceSuite) TestValidateToken() {
	s.Run("Success", func() {
		user := domain.User{
			ID:       "1",
			Username: "Abebe",
			Role:     "user",
		}
		tokens, err := s.jwtService.GenerateTokens(user)
		s.NoError(err)

		claims, err := s.jwtService.ValidateToken(tokens.AccessToken)

		s.NoError(err)
		s.Equal(user.ID, claims["sub"])
		s.Equal(user.Username, claims["username"])
		s.Equal(user.Role, claims["role"])
	})

	s.Run("InvalidToken", func() {
		claims, err := s.jwtService.ValidateToken("invalid_token")

		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("WrongSecret", func() {
		// Create a token with a different secret
		wrongService := security.NewJWTService("wrong_secret", "refresh_secret", 1, 24).(*security.JwtService)
		user := domain.User{ID: "1", Username: "Abebe", Role: "user"}
		tokens, err := wrongService.GenerateTokens(user)
		s.NoError(err)
		claims, err := s.jwtService.ValidateToken(tokens.AccessToken)
		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("ExpiredToken", func() {
		// Create a service with a very short expiry
		shortExpiryService := security.NewJWTService("access_secret", "refresh_secret", 0, 24).(*security.JwtService)
		user := domain.User{ID: "1", Username: "Abebe", Role: "user"}
		tokens, err := shortExpiryService.GenerateTokens(user)
		s.NoError(err)

		// Wait for token to expire (simulate by setting expiry to 0 hours)
		claims, err := shortExpiryService.ValidateToken(tokens.AccessToken)

		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("WrongSigningMethod", func() {
		// Create a token with a different signing method (e.g., HS512)
		claims := jwt.MapClaims{
			"sub":      "1",
			"username": "Abebe",
			"role":     "user",
			"exp":      time.Now().Add(time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		tokenStr, err := token.SignedString([]byte("access_secret"))
		s.NoError(err)

		claimsOut, err := s.jwtService.ValidateToken(tokenStr)

		s.Error(err)
		s.Contains(err.Error(), "unexpected signing method")
		s.Nil(claimsOut)
	})
}

// TestValidateRefreshToken tests the ValidateRefreshToken function
func (s *JwtServiceSuite) TestValidateRefreshToken() {
	s.Run("Success", func() {
		user := domain.User{
			ID:       "1",
			Username: "Abebe",
			Role:     "user",
		}
		tokens, err := s.jwtService.GenerateTokens(user)
		s.NoError(err)

		claims, err := s.jwtService.ValidateRefreshToken(tokens.RefreshToken)

		s.NoError(err)
		s.Equal(user.ID, claims["sub"])
		s.Equal(user.Username, claims["username"])
		s.NotContains(claims, "role") // Refresh token does not include role
	})

	s.Run("InvalidToken", func() {
		claims, err := s.jwtService.ValidateRefreshToken("invalid_token")

		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("WrongSecret", func() {
		// Create a token with a different refresh secret
		wrongService := security.NewJWTService("access_secret", "wrong_secret", 1, 24).(*security.JwtService)
		user := domain.User{ID: "1", Username: "Abebe", Role: "user"}
		tokens, err := wrongService.GenerateTokens(user)
		s.NoError(err)

		claims, err := s.jwtService.ValidateRefreshToken(tokens.RefreshToken)

		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("ExpiredToken", func() {
		// Create a service with a very short refresh expiry
		shortExpiryService := security.NewJWTService("access_secret", "refresh_secret", 1, 0).(*security.JwtService)
		user := domain.User{ID: "1", Username: "Abebe", Role: "user"}
		tokens, err := shortExpiryService.GenerateTokens(user)
		s.NoError(err)

		// Wait for token to expire (simulate by setting expiry to 0 hours)
		claims, err := shortExpiryService.ValidateRefreshToken(tokens.RefreshToken)

		s.Error(err)
		s.Contains(err.Error(), "invalid token")
		s.Nil(claims)
	})

	s.Run("WrongSigningMethod", func() {
		// Create a refresh token with a different signing method (e.g., HS512)
		claims := jwt.MapClaims{
			"sub":      "1",
			"username": "Abebe",
			"exp":      time.Now().Add(time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		tokenStr, err := token.SignedString([]byte("refresh_secret"))
		s.NoError(err)

		claimsOut, err := s.jwtService.ValidateRefreshToken(tokenStr)

		s.Error(err)
		s.Contains(err.Error(), "unexpected signing method")
		s.Nil(claimsOut)
	})
}
