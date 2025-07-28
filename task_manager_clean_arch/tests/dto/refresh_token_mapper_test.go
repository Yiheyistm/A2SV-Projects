package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
)

// RefreshTokenMapperSuite defines the test suite for refresh token mapper functions
type RefreshTokenMapperSuite struct {
	suite.Suite
}

// TestRefreshTokenMapperSuite runs the test suite
func TestRefreshTokenMapperSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenMapperSuite))
}

// TestFromRequestToDomainRefreshToken tests the FromRequestToDomainRefreshToken function
func (s *RefreshTokenMapperSuite) TestFromRequestToDomainRefreshToken() {
	s.Run("Success", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{
			RefreshToken: "valid_refresh_token",
		}
		expectedRefreshToken := &domain.RefreshToken{
			RefreshToken: "valid_refresh_token",
		}

		result := refreshTokenRequest.FromRequestToDomainRefreshToken()

		s.Equal(expectedRefreshToken, result)
	})

	s.Run("EmptyField", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{
			RefreshToken: "",
		}
		expectedRefreshToken := &domain.RefreshToken{
			RefreshToken: "",
		}

		result := refreshTokenRequest.FromRequestToDomainRefreshToken()

		s.Equal(expectedRefreshToken, result)
	})
}

// TestFromDomainRefreshTokenToResponse tests the FromDomainRefreshTokenToResponse function
func (s *RefreshTokenMapperSuite) TestFromDomainRefreshTokenToResponse() {
	s.Run("Success", func() {
		domainRefreshToken := &domain.RefreshToken{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}
		expectedResponse := &domain.RefreshToken{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}

		result := dto.FromDomainRefreshTokenToResponse(domainRefreshToken)

		s.Equal(expectedResponse, result)
	})

	s.Run("EmptyFields", func() {
		domainRefreshToken := &domain.RefreshToken{
			AccessToken:  "",
			RefreshToken: "",
		}
		expectedResponse := &domain.RefreshToken{
			AccessToken:  "",
			RefreshToken: "",
		}

		result := dto.FromDomainRefreshTokenToResponse(domainRefreshToken)

		s.Equal(expectedResponse, result)
	})

	s.Run("NilInput", func() {
		var domainRefreshToken *domain.RefreshToken
		result := dto.FromDomainRefreshTokenToResponse(domainRefreshToken)

		s.Nil(result)
	})
}
