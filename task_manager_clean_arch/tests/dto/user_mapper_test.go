package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
)

// UserMapperSuite defines the test suite for user_mapper.go
type UserMapperSuite struct {
	suite.Suite
}

// TestUserMapperSuite runs the test suite
func TestUserMapperSuite(t *testing.T) {
	suite.Run(t, new(UserMapperSuite))
}

// TestFromRequestToDomainUser tests the FromRequestToDomainUser function
func (s *UserMapperSuite) TestFromRequestToDomainUser() {
	s.Run("Success", func() {
		userRequest := dto.UserRequest{
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}
		expectedUser := &domain.User{
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}

		result := userRequest.FromRequestToDomainUser()

		s.Equal(expectedUser, result)
	})

	s.Run("EmptyFields", func() {
		userRequest := dto.UserRequest{
			Username: "",
			Email:    "",
			Password: "",
			Role:     "",
		}
		expectedUser := &domain.User{
			Username: "",
			Email:    "",
			Password: "",
			Role:     "",
		}

		result := userRequest.FromRequestToDomainUser()

		s.Equal(expectedUser, result)
	})
}

// TestFromDomainUserToResponse tests the FromDomainUserToResponse function
func (s *UserMapperSuite) TestFromDomainUserToResponse() {
	s.Run("Success", func() {
		domainUser := &domain.User{
			ID:       "1",
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}
		expectedResponse := &dto.UserResponse{
			ID:       "1",
			Username: "Abebe",
			Email:    "abebe@example.com",
			Role:     "user",
		}

		result := dto.FromDomainUserToResponse(domainUser)

		s.Equal(expectedResponse, result)
	})

	s.Run("EmptyFields", func() {
		domainUser := &domain.User{
			ID:       "",
			Username: "",
			Email:    "",
			Password: "",
			Role:     "",
		}
		expectedResponse := &dto.UserResponse{
			ID:       "",
			Username: "",
			Email:    "",
			Role:     "",
		}

		result := dto.FromDomainUserToResponse(domainUser)

		s.Equal(expectedResponse, result)
	})

	s.Run("NilInput", func() {
		var domainUsers domain.User
		result := dto.FromDomainUserToResponse(&domainUsers)

		s.Equal(&dto.UserResponse{}, result)
	})
}

// TestFromDomainUserToResponseList tests the FromDomainUserToResponseList function
func (s *UserMapperSuite) TestFromDomainUserToResponseList() {
	s.Run("Success", func() {
		domainUsers := []domain.User{
			{
				ID:       "1",
				Username: "Abebe",
				Email:    "abebe@example.com",
				Password: "password123",
				Role:     "user",
			},
			{
				ID:       "2",
				Username: "Kebede",
				Email:    "kebede@example.com",
				Password: "password456",
				Role:     "admin",
			},
		}
		expectedResponses := []dto.UserResponse{
			{
				ID:       "1",
				Username: "Abebe",
				Email:    "abebe@example.com",
				Role:     "user",
			},
			{
				ID:       "2",
				Username: "Kebede",
				Email:    "kebede@example.com",
				Role:     "admin",
			},
		}

		result := dto.FromDomainUserToResponseList(domainUsers)

		s.Equal(expectedResponses, result)
	})

	s.Run("EmptySlice", func() {
		domainUsers := []domain.User{}
		var expectedResponses []dto.UserResponse

		result := dto.FromDomainUserToResponseList(domainUsers)

		s.Equal(expectedResponses, result)
	})

	s.Run("NilSlice", func() {
		var domainUsers []domain.User
		var expectedResponses []dto.UserResponse

		result := dto.FromDomainUserToResponseList(domainUsers)

		s.Equal(expectedResponses, result)
	})
}
