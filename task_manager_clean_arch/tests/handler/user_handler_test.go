package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/handler"
	"github.com/yiheyistm/task_manager/mocks/mocks_domain"
	"github.com/yiheyistm/task_manager/mocks/mocks_security"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserHandlerSuite defines the test suite for UserHandler
type UserHandlerSuite struct {
	suite.Suite
	mockUserUsecase         *mocks_domain.IUserUseCase
	mockTaskUsecase         *mocks_domain.ITaskUseCase
	mockRefreshTokenUsecase *mocks_security.IRefreshTokenUsecase
	handler                 *handler.UserHandler
	validate                *validator.Validate
}

// SetupTest initializes the mocks and handler before each test
func (s *UserHandlerSuite) SetupTest() {
	s.mockUserUsecase = mocks_domain.NewIUserUseCase(s.T())
	s.mockTaskUsecase = mocks_domain.NewITaskUseCase(s.T())
	s.mockRefreshTokenUsecase = mocks_security.NewIRefreshTokenUsecase(s.T())
	s.handler = &handler.UserHandler{
		UserUsecase:         s.mockUserUsecase,
		TaskUsecase:         s.mockTaskUsecase,
		RefreshTokenUsecase: s.mockRefreshTokenUsecase,
	}
	s.validate = validator.New()
	// validate := s.validate
}

// TestUserHandlerSuite runs the test suite
func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

// TestRegisterRequest tests the RegisterRequest method
func (s *UserHandlerSuite) TestRegisterRequest() {
	s.Run("Success", func() {
		userRequest := dto.UserRequest{
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}
		user := domain.User{
			Username: strings.ToLower(userRequest.Username),
			Email:    strings.ToLower(userRequest.Email),
			Password: "hashed_password",
			Role:     userRequest.Role,
		}
		s.mockUserUsecase.On("GetByUsername", strings.ToLower(userRequest.Username)).Return(nil, nil)
		s.mockUserUsecase.On("Insert", mock.MatchedBy(func(u *domain.User) bool {
			return u.Username == user.Username &&
				u.Email == user.Email &&
				u.Role == user.Role &&
				len(u.Password) == 60
		})).Return(nil)

		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RegisterRequest(c)

		s.Equal(http.StatusCreated, w.Code)
		var response dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(strings.ToLower(userRequest.Username), response.Username)
		s.Equal(userRequest.Email, response.Email)
		s.Equal(userRequest.Role, response.Role)
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RegisterRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		userRequest := dto.UserRequest{
			Username: "ab",      // Too short
			Email:    "invalid", // Invalid email
			Password: "pass",    // Too short
			Role:     "invalid", // Invalid role
		}
		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RegisterRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
		s.resetMocks()
	})

	s.Run("UsernameExists", func() {
		s.mockUserUsecase.ExpectedCalls = nil
		userRequest := dto.UserRequest{
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}
		existingUser := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetByUsername", strings.ToLower(userRequest.Username)).Return(existingUser, nil)

		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RegisterRequest(c)
		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Username already exists", response["error"])
		s.mockUserUsecase.AssertCalled(s.T(), "GetByUsername", strings.ToLower(userRequest.Username))
		s.resetMocks()
	})

	s.Run("InsertError", func() {
		s.mockUserUsecase.ExpectedCalls = nil
		userRequest := dto.UserRequest{
			Username: "Abebe",
			Email:    "abebe@example.com",
			Password: "password123",
			Role:     "user",
		}
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := domain.User{
			Username: strings.ToLower(userRequest.Username),
			Email:    strings.ToLower(userRequest.Email),
			Password: string(hashed_password),
			Role:     userRequest.Role,
		}
		s.mockUserUsecase.On("GetByUsername", strings.ToLower(userRequest.Username)).Return(nil, nil)
		s.mockUserUsecase.On("Insert", mock.MatchedBy(func(u *domain.User) bool {
			return u.Username == user.Username &&
				u.Email == user.Email &&
				u.Role == user.Role &&
				len(u.Password) == 60
		})).Return(errors.New("insert failed"))

		body, _ := json.Marshal(userRequest)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RegisterRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to register user", response["error"])
		s.resetMocks()
	})
}

// TestLoginRequest tests the LoginRequest method
func (s *UserHandlerSuite) TestLoginRequest() {
	s.Run("SuccessEmail", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "abebe@example.com",
			Password:   "password123",
		}
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       "1",
			Username: "abebe",
			Email:    "abebe@example.com",
			Password: string(hashed_password),
			Role:     "user",
		}
		tokens := domain.RefreshToken{AccessToken: "access_token", RefreshToken: "refresh_token"}
		s.mockUserUsecase.On("GetByEmail", loginRequest.Identifier).Return(user, nil)
		s.mockRefreshTokenUsecase.On("GenerateTokens", *user).Return(tokens, nil)

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(tokens.AccessToken, response.AccessToken)
		s.Equal(tokens.RefreshToken, response.RefreshToken)
		s.resetMocks()
	})

	s.Run("SuccessUsername", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "abebe",
			Password:   "password123",
		}
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       "1",
			Username: "abebe",
			Email:    "abebe@example.com",
			Password: string(hashed_password),
			Role:     "user",
		}
		tokens := domain.RefreshToken{AccessToken: "access_token", RefreshToken: "refresh_token"}
		s.mockUserUsecase.On("GetByUsername", loginRequest.Identifier).Return(user, nil)
		s.mockRefreshTokenUsecase.On("GenerateTokens", *user).Return(tokens, nil)

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(tokens.AccessToken, response.AccessToken)
		s.Equal(tokens.RefreshToken, response.RefreshToken)
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("EmptyIdentifierOrPassword", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "",
			Password:   "",
		}
		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Username/email and password are required", response["message"])
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		s.mockUserUsecase.ExpectedCalls = nil
		loginRequest := dto.LoginRequest{
			Identifier: "abebe@example.com",
			Password:   "password123",
		}
		s.mockUserUsecase.On("GetByEmail", loginRequest.Identifier).Return(nil, errors.New("user not found"))

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid credentials", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidPassword", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "abebe@example.com",
			Password:   "wrongpassword",
		}
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       "1",
			Username: "abebe",
			Email:    "abebe@example.com",
			Password: string(hashed_password),
			Role:     "user",
		}
		s.mockUserUsecase.On("GetByEmail", loginRequest.Identifier).Return(user, nil)

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid email or password", response["error"])
		s.resetMocks()
	})

	s.Run("TokenGenerationError", func() {
		s.mockUserUsecase.ExpectedCalls = nil
		s.mockRefreshTokenUsecase.ExpectedCalls = nil
		loginRequest := dto.LoginRequest{
			Identifier: "abebe@example.com",
			Password:   "password123",
		}
		hashed_password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       "1",
			Username: "abebe",
			Email:    "abebe@example.com",
			Password: string(hashed_password),
			Role:     "user",
		}
		s.mockUserUsecase.On("GetByEmail", loginRequest.Identifier).Return(user, nil)
		s.mockRefreshTokenUsecase.On("GenerateTokens", *user).Return(domain.RefreshToken{}, errors.New("token generation failed"))

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.LoginRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to generate token", response["error"])
		s.resetMocks()
	})
}

// TestGetAllUsers tests the GetAllUsers method
func (s *UserHandlerSuite) TestGetAllUsers() {
	s.Run("Success", func() {
		users := []domain.User{
			{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"},
			{ID: "2", Username: "kebede", Email: "kebede@example.com", Role: "admin"},
		}
		s.mockUserUsecase.On("GetAll").Return(users, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.GetAllUsers(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Len(response["users"], 2)
		s.resetMocks()
	})

	s.Run("FetchError", func() {
		s.mockUserUsecase.On("GetAll").Return(nil, errors.New("fetch failed"))

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.GetAllUsers(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to fetch users", response["error"])
		s.resetMocks()
	})
}

// TestGetUser tests the GetUser method
func (s *UserHandlerSuite) TestGetUser() {
	s.Run("Success", func() {
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockUserUsecase.On("GetByUsername", "abebe").Return(user, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUser(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("abebe", response["user"].(map[string]interface{})["username"])
		s.resetMocks()
	})

	s.Run("EmptyUsername", func() {
		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: ""}}

		s.handler.GetUser(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("User name is required", response["error"])
		s.resetMocks()
	})

	s.Run("FetchError", func() {
		s.mockUserUsecase.On("GetByUsername", "abebe").Return(nil, errors.New("user not found"))

		req := httptest.NewRequest(http.MethodGet, "/users/abebe", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUser(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to fetch user", response["error"])
		s.resetMocks()
	})
}

// TestGetUserTasks tests the GetUserTasks method
func (s *UserHandlerSuite) TestGetUserTasks() {
	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()},
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetTasksByUser", "abebe").Return(tasks, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTasks(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Len(response["tasks"], 1)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTasks(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to see details about this user", response["error"])
		s.resetMocks()
	})

	s.Run("FetchError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetTasksByUser", "abebe").Return(nil, errors.New("fetch failed"))

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTasks(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to fetch tasks for user", response["error"])
		s.resetMocks()
	})
}

// TestGetUserTask tests the GetUserTask method
func (s *UserHandlerSuite) TestGetUserTask() {
	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		id := primitive.NewObjectID()
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: dueDate}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetByIdAndUser", id.Hex(), "abebe").Return(task, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/"+id.Hex(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: id.Hex()}}

		s.handler.GetUserTask(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(task.Title, response.Title)
		s.Equal(task.Description, response.Description)
		s.Equal(task.Status, response.Status)
		s.Equal(task.CreatedBy, response.CreatedBy)
		s.Equal(task.DueDate, response.DueDate)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.GetUserTask(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to see details about this user", response["error"])
		s.resetMocks()
	})

	s.Run("EmptyTaskID", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: ""}}

		s.handler.GetUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
		s.resetMocks()
	})

	s.Run("TaskNotFound", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetByIdAndUser", "1", "abebe").Return(domain.Task{}, errors.New("task not found"))

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.GetUserTask(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("task not found", response["message"])
		s.resetMocks()
	})
}

// TestCreateUserTask tests the CreateUserTask method
func (s *UserHandlerSuite) TestCreateUserTask() {

	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		request := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "pending",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Create", mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == request.Title &&
				t.Description == request.Description &&
				t.Status == request.Status &&
				t.CreatedBy == user.Username &&
				t.DueDate.Equal(dueDate)
		})).Return(nil)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/users/abebe/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}
		s.handler.CreateUserTask(c)

		s.Equal(http.StatusCreated, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(request.Title, response.Title)
		s.Equal(request.CreatedBy, response.CreatedBy)
		s.Equal(request.Description, response.Description)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		task := dto.TaskRequest{Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending"}
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPost, "/users/abebe/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.CreateUserTask(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to create tasks on behalf of other user", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodPost, "/users/abebe/tasks", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.CreateUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		s.resetMocks()
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		task := domain.Task{Title: "", Status: "", CreatedBy: "abebe"} // Invalid fields
		s.mockTaskUsecase.On("Create", &task).Return(errors.New("Field validation"))
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPost, "/users/abebe/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.CreateUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
		s.resetMocks()
	})

	s.Run("CreateError", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)

		task := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "pending",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Create", mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == task.Title &&
				t.Description == task.Description &&
				t.Status == task.Status &&
				t.CreatedBy == task.CreatedBy
		})).Return(errors.New("Failed to create task"))

		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPost, "/users/abebe/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.CreateUserTask(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to create task", response["error"])
		s.resetMocks()
	})
}

// TestUpdateUserTask tests the UpdateUserTask method
func (s *UserHandlerSuite) TestUpdateUserTask() {

	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		id := primitive.NewObjectID()
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "completed",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("UpdateByIdAndUser", id.Hex(), mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == task.Title &&
				t.Description == task.Description &&
				t.Status == task.Status &&
				t.CreatedBy == task.CreatedBy &&
				t.DueDate.Equal(task.DueDate)
		}), "abebe").Return(nil)

		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: id.Hex()}}

		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(task.CreatedBy, response.CreatedBy)
		s.Equal(task.Title, response.Title)
		s.Equal(task.Description, response.Description)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		task := dto.TaskRequest{Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed"}
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to update this task", response["error"])
		s.resetMocks()
	})

	s.Run("EmptyTaskID", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		task := dto.TaskRequest{Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed"}
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: ""}}

		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/1", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		task := domain.Task{Title: "", Status: ""} // Invalid fields
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}
		fmt.Println("Running validation error test", c)
		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
		s.resetMocks()
	})

	s.Run("UpdateError", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		task := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "completed",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("UpdateByIdAndUser", "1", mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == task.Title &&
				t.Description == task.Description &&
				t.Status == task.Status &&
				t.CreatedBy == task.CreatedBy &&
				t.DueDate.Equal(task.DueDate)
		}), "abebe").Return(errors.New("Failed to update task"))
		body, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/users/abebe/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.UpdateUserTask(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to update task", response["error"])
		s.resetMocks()
	})
}

// TestDeleteUserTask tests the DeleteUserTask method
func (s *UserHandlerSuite) TestDeleteUserTask() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("DeleteByIdAndUser", id.Hex(), "abebe").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/users/abebe/tasks/"+id.Hex(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: id.Hex()}}

		s.handler.DeleteUserTask(c)

		s.Equal(http.StatusNoContent, w.Code)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodDelete, "/users/abebe/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.DeleteUserTask(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to delete tasks on behalf of other user", response["error"])
		s.resetMocks()
	})

	s.Run("EmptyTaskID", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodDelete, "/users/abebe/tasks/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: ""}}

		s.handler.DeleteUserTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
		s.resetMocks()
	})

	s.Run("DeleteError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("DeleteByIdAndUser", "1", "abebe").Return(errors.New("delete failed"))

		req := httptest.NewRequest(http.MethodDelete, "/users/abebe/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}, {Key: "id", Value: "1"}}

		s.handler.DeleteUserTask(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to delete task", response["error"])
		s.resetMocks()
	})
}

// TestGetUserTaskStats tests the GetUserTaskStats method
func (s *UserHandlerSuite) TestGetUserTaskStats() {
	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		stats := []domain.StatusCount{
			{Status: "pending", Count: 5},
			{Status: "completed", Count: 3},
		}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetTaskStatsByUser", "abebe").Return(stats, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/stats", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTaskStats(c)

		s.Equal(http.StatusOK, w.Code)
		var response []domain.StatusCount
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(stats, response)
		s.resetMocks()
	})

	s.Run("PermissionDenied", func() {
		user := &domain.User{Username: "kebede"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/stats", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTaskStats(c)

		s.Equal(http.StatusForbidden, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You do not have permission to see details about this user", response["error"])
		s.resetMocks()
	})

	s.Run("FetchError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("GetTaskStatsByUser", "abebe").Return(nil, errors.New("stats fetch failed"))

		req := httptest.NewRequest(http.MethodGet, "/users/abebe/tasks/stats", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "abebe"}}

		s.handler.GetUserTaskStats(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to fetch task stats", response["error"])
		s.resetMocks()
	})
}

func (s *UserHandlerSuite) resetMocks() {
	s.mockUserUsecase.ExpectedCalls = nil
	s.mockUserUsecase.Calls = nil
	s.mockTaskUsecase.ExpectedCalls = nil
	s.mockTaskUsecase.Calls = nil
	s.mockRefreshTokenUsecase.ExpectedCalls = nil
	s.mockRefreshTokenUsecase.Calls = nil
}
