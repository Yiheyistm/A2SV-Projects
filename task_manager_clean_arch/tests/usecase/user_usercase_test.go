package usecase

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/usecase"
	"github.com/yiheyistm/task_manager/mocks/mocks_domain"
)

// UserUseCaseSuite defines the test suite for UserUseCase
type UserUseCaseSuite struct {
	suite.Suite
	mockRepo *mocks_domain.UserRepository
	useCase  domain.IUserUseCase
}

// SetupTest initializes the mocks and use case before each test
func (s *UserUseCaseSuite) SetupTest() {
	s.mockRepo = mocks_domain.NewUserRepository(s.T())
	s.useCase = usecase.NewUserUseCase(s.mockRepo)
}

// TestUserUseCaseSuite runs the test suite
func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

// TestGetAll tests the GetAll method
func (s *UserUseCaseSuite) TestGetAll() {
	s.Run("Success", func() {
		users := []domain.User{
			{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"},
			{ID: "2", Username: "kebede", Email: "kebede@example.com", Role: "admin"},
		}
		s.mockRepo.On("GetAll", mock.Anything).Return(users, nil)

		result, err := s.useCase.GetAll()

		s.NoError(err)
		s.Equal(users, result)
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		s.mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))

		result, err := s.useCase.GetAll()

		s.Error(err)
		s.EqualError(err, "database error")
		s.Nil(result)
	})
}

// // TestGetByUsername tests the GetByUsername method
func (s *UserUseCaseSuite) TestGetByUsername() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockRepo.On("GetByUsername", mock.Anything, "abebe").Return(user, nil)

		result, err := s.useCase.GetByUsername("abebe")

		s.NoError(err)
		s.Equal(user, result)
	})

	s.Run("UserNotFound", func() {
		s.mockRepo.On("GetByUsername", mock.Anything, "unknown").Return(nil, errors.New("user not found"))

		result, err := s.useCase.GetByUsername("unknown")

		s.Error(err)
		s.EqualError(err, "user not found")
		s.Nil(result)
	})
}

// TestGetByEmail tests the GetByEmail method
func (s *UserUseCaseSuite) TestGetByEmail() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockRepo.On("GetByEmail", mock.Anything, "abebe@example.com").Return(user, nil)

		result, err := s.useCase.GetByEmail("abebe@example.com")

		s.NoError(err)
		s.Equal(user, result)
	})

	s.Run("EmailNotFound", func() {
		s.mockRepo.On("GetByEmail", mock.Anything, "unknown@example.com").Return(nil, errors.New("email not found"))

		result, err := s.useCase.GetByEmail("unknown@example.com")

		s.Error(err)
		s.EqualError(err, "email not found")
		s.Nil(result)
	})
}

// TestInsert tests the Insert method
func (s *UserUseCaseSuite) TestInsert() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockRepo.On("Insert", mock.Anything, user).Return(nil)

		err := s.useCase.Insert(user)

		s.NoError(err)
	})

	s.Run("NilUser", func() {
		err := s.useCase.Insert(nil)

		s.Error(err)
		s.EqualError(err, "user cannot be nil")
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockRepo.On("Insert", mock.Anything, user).Return(errors.New("insert failed"))

		err := s.useCase.Insert(user)

		s.Error(err)
		s.EqualError(err, "insert failed")
	})
}

// TestGetUserFromContext tests the GetUserFromContext method
func (s *UserUseCaseSuite) TestGetUserFromContext() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		ginContext := &gin.Context{}
		s.mockRepo.On("GetUserFromContext", ginContext).Return(user)

		result := s.useCase.GetUserFromContext(ginContext)

		s.Equal(user, result)
	})

	s.Run("NoUserInContext", func() {
		s.mockRepo.ExpectedCalls = nil
		ginContext := &gin.Context{}
		s.mockRepo.On("GetUserFromContext", ginContext).Return(nil)

		result := s.useCase.GetUserFromContext(ginContext)

		s.Equal(&domain.User{}, result)
	})
}
