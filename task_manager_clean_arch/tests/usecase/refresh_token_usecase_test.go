package usecase

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/usecase"
	"github.com/yiheyistm/task_manager/mocks/mocks_domain"
	"github.com/yiheyistm/task_manager/mocks/mocks_security"
)

// RefreshTokenUsecaseSuite defines the test suite for refreshTokenUsecase
type RefreshTokenUsecaseSuite struct {
	suite.Suite
	mockUserRepo *mocks_domain.UserRepository
	mockJwt      *mocks_security.IRefreshTokenUsecase
	useCase      domain.IRefreshTokenUsecase
}

// SetupTest initializes the mocks and use case before each test
func (s *RefreshTokenUsecaseSuite) SetupTest() {
	s.mockUserRepo = mocks_domain.NewUserRepository(s.T())
	s.mockJwt = mocks_security.NewIRefreshTokenUsecase(s.T())
	s.useCase = usecase.NewRefreshTokenUsecase(s.mockUserRepo, s.mockJwt)
}

// TestRefreshTokenUsecaseSuite runs the test suite
func TestRefreshTokenUsecaseSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenUsecaseSuite))
}

// TestGetByUsername tests the GetByUsername method
func (s *RefreshTokenUsecaseSuite) TestGetByUsername() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockUserRepo.On("GetByUsername", mock.Anything, "abebe").Return(user, nil)

		result, err := s.useCase.GetByUsername("abebe")

		s.NoError(err)
		s.Equal(user, result)
	})

	s.Run("UserNotFound", func() {
		s.mockUserRepo.On("GetByUsername", mock.Anything, "unknown").Return(nil, errors.New("user not found"))

		result, err := s.useCase.GetByUsername("unknown")

		s.Error(err)
		s.EqualError(err, "user not found")
		s.Nil(result)
	})
}

// TestGenerateTokens tests the GenerateTokens method
func (s *RefreshTokenUsecaseSuite) TestGenerateTokens() {
	s.Run("Success", func() {
		user := domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		expectedTokens := domain.RefreshToken{AccessToken: "access_token", RefreshToken: "refresh_token"}
		s.mockJwt.On("GenerateTokens", user).Return(expectedTokens, nil)

		result, err := s.useCase.GenerateTokens(user)
		s.NoError(err)
		s.Equal(expectedTokens.AccessToken, result.AccessToken)
		s.Equal(expectedTokens.RefreshToken, result.RefreshToken)
		s.mockJwt.AssertExpectations(s.T())
	})
	s.Run("JwtError", func() {
		s.mockJwt.ExpectedCalls = nil
		user := domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockJwt.On("GenerateTokens", user).Return(domain.RefreshToken{}, errors.New("jwt generation failed"))
		result, err := s.useCase.GenerateTokens(user)
		s.Error(err)
		s.EqualError(err, "jwt generation failed")
		s.Equal(domain.RefreshToken{}, result)

		s.mockJwt.AssertExpectations(s.T())
	})
}

// TestValidateRefreshToken tests the ValidateRefreshToken method
func (s *RefreshTokenUsecaseSuite) TestValidateRefreshToken() {
	s.Run("Success", func() {
		claims := jwt.MapClaims{"sub": "1", "username": "abebe"}
		s.mockJwt.On("ValidateRefreshToken", "valid_refresh_token").Return(claims, nil)

		result, err := s.useCase.ValidateRefreshToken("valid_refresh_token")

		s.NoError(err)
		s.Equal(claims, result)
	})

	s.Run("InvalidToken", func() {
		s.mockJwt.On("ValidateRefreshToken", "invalid_refresh_token").Return(jwt.MapClaims{}, errors.New("invalid refresh token"))

		result, err := s.useCase.ValidateRefreshToken("invalid_refresh_token")

		s.Error(err)
		s.EqualError(err, "invalid refresh token")
		s.Equal(jwt.MapClaims{}, result)
	})
}

// TestValidateToken tests the ValidateToken method
func (s *RefreshTokenUsecaseSuite) TestValidateToken() {
	s.Run("Success", func() {
		claims := jwt.MapClaims{"sub": "1", "username": "abebe"}
		s.mockJwt.On("ValidateToken", "valid_access_token").Return(claims, nil)

		result, err := s.useCase.ValidateToken("valid_access_token")

		s.NoError(err)
		s.Equal(claims, result)
	})

	s.Run("InvalidToken", func() {
		s.mockJwt.On("ValidateToken", "invalid_access_token").Return(jwt.MapClaims{}, errors.New("invalid access token"))

		result, err := s.useCase.ValidateToken("invalid_access_token")

		s.Error(err)
		s.EqualError(err, "invalid access token")
		s.Equal(jwt.MapClaims{}, result)
	})
}
