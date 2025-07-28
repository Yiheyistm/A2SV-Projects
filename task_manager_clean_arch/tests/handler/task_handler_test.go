package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
	mocks_domain "github.com/yiheyistm/task_manager/mocks/mocks_domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskHandlerSuite defines the test suite for TaskHandler
type TaskHandlerSuite struct {
	suite.Suite
	mockTaskUsecase *mocks_domain.ITaskUseCase
	mockUserUsecase *mocks_domain.IUserUseCase
	handler         *handler.TaskHandler
	validate        *validator.Validate
}

// SetupTest initializes the mocks and handler before each test
func (s *TaskHandlerSuite) SetupTest() {
	s.mockTaskUsecase = mocks_domain.NewITaskUseCase(s.T())
	s.mockUserUsecase = mocks_domain.NewIUserUseCase(s.T())
	s.handler = &handler.TaskHandler{
		TaskUsecase: s.mockTaskUsecase,
		UserUsecase: s.mockUserUsecase,
	}
	s.validate = validator.New()
}

// TestTaskHandlerSuite runs the test suite
func TestTaskHandlerSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlerSuite))
}

// TestGetTasks tests the GetTasks method
func (s *TaskHandlerSuite) TestGetTasks() {
	s.Run("Success", func() {
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: dueDate},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", Description: "Trade in Merkato", Status: "completed", CreatedBy: "kebede", DueDate: dueDate},
		}
		s.mockTaskUsecase.On("GetAll").Return(tasks, nil)
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		s.handler.GetTasks(c)

		s.Equal(http.StatusOK, w.Code)
		var response []dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Len(response, 2)
		s.Equal(tasks[0].Title, response[0].Title)
		s.Equal(tasks[1].Title, response[1].Title)
	})
	s.resetMocks()

	s.Run("FetchError", func() {
		s.mockTaskUsecase.On("GetAll").Return(nil, errors.New("fetch failed"))

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.GetTasks(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to retrieve tasks", response["error"])
	})
	s.resetMocks()
}

// TestGetTask tests the GetTask method
func (s *TaskHandlerSuite) TestGetTask() {
	s.Run("Success", func() {
		task := domain.Task{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe"}
		s.mockTaskUsecase.On("GetById", "1").Return(task, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.GetTask(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(task.Title, response.Title)
	})
	s.resetMocks()

	s.Run("EmptyID", func() {
		req := httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: ""}}

		s.handler.GetTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
	})
	s.resetMocks()

	s.Run("TaskNotFound", func() {
		s.mockTaskUsecase.On("GetById", "1").Return(domain.Task{}, errors.New("task not found"))

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.GetTask(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("task not found", response["message"])
	})
	s.resetMocks()
}

// TestCreateTask tests the CreateTask method
func (s *TaskHandlerSuite) TestCreateTask() {
	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		taskRequest := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "pending",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		task := taskRequest.FromRequestToDomainTask()
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Create", task).Return(nil)

		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.CreateTask(c)

		s.Equal(http.StatusCreated, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(taskRequest.Title, response.Title)
		s.Equal(taskRequest.CreatedBy, response.CreatedBy)
	})
	s.resetMocks()
	s.Run("InvalidJSON", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.CreateTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
	})
	s.resetMocks()
	s.Run("ValidationError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		taskRequest := dto.TaskRequest{Title: "", Status: ""} // Invalid fields
		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		req.Header.Set("ContentType", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.CreateTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
	})
	s.resetMocks()
	s.Run("CreateError", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		taskRequest := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "pending",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		task := taskRequest.FromRequestToDomainTask()
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Create", task).Return(errors.New("create failed"))

		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.CreateTask(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to create task", response["error"])
	})
}

// TestUpdateTask tests the UpdateTask method
func (s *TaskHandlerSuite) TestUpdateTask() {
	s.Run("Success", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		taskRequest := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "completed",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		task := taskRequest.FromRequestToDomainTask()
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Update", "1", task).Return(nil)

		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.UpdateTask(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.TaskResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(taskRequest.Title, response.Title)
		s.Equal(taskRequest.CreatedBy, response.CreatedBy)
	})
	s.resetMocks()

	s.Run("EmptyID", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		taskRequest := dto.TaskRequest{Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed"}
		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPut, "/tasks/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: ""}}

		s.handler.UpdateTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
	})
	s.resetMocks()

	s.Run("InvalidJSON", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		req := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.UpdateTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
	})
	s.resetMocks()

	s.Run("ValidationError", func() {
		user := &domain.User{Username: "abebe"}
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)

		taskRequest := dto.TaskRequest{Title: "", Status: ""} // Invalid fields
		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.UpdateTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
	})
	s.resetMocks()

	s.Run("UpdateError", func() {
		user := &domain.User{Username: "abebe"}
		dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		taskRequest := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			Status:      "completed",
			CreatedBy:   "abebe",
			DueDate:     dueDate,
		}
		task := taskRequest.FromRequestToDomainTask()
		s.mockUserUsecase.On("GetUserFromContext", mock.Anything).Return(user)
		s.mockTaskUsecase.On("Update", "1", task).Return(errors.New("update failed"))

		body, _ := json.Marshal(taskRequest)
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.UpdateTask(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to update task", response["message"])
	})
}

// TestDeleteTask tests the DeleteTask method
func (s *TaskHandlerSuite) TestDeleteTask() {
	s.Run("Success", func() {
		s.mockTaskUsecase.On("Delete", "1").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.DeleteTask(c)

		s.Equal(http.StatusNoContent, w.Code)
	})
	s.resetMocks()

	s.Run("EmptyID", func() {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: ""}}

		s.handler.DeleteTask(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Task ID is required", response["error"])
	})
	s.resetMocks()

	s.Run("DeleteError", func() {
		s.mockTaskUsecase.On("Delete", "1").Return(errors.New("delete failed"))

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		s.handler.DeleteTask(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("delete failed", response["message"])
	})
}

// TestGetTaskCountByStatus tests the GetTaskCountByStatus method
func (s *TaskHandlerSuite) TestGetTaskCountByStatus() {
	s.Run("Success", func() {
		counts := []domain.StatusCount{
			{Status: "pending", Count: 5},
			{Status: "completed", Count: 3},
		}
		s.mockTaskUsecase.On("GetTaskCountByStatus").Return(counts, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks/status", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.GetTaskCountByStatus(c)

		s.Equal(http.StatusOK, w.Code)
		var response []domain.StatusCount
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(counts, response)
	})
	s.resetMocks()
	s.Run("FetchError", func() {
		s.mockTaskUsecase.On("GetTaskCountByStatus").Return(nil, errors.New("fetch failed"))

		req := httptest.NewRequest(http.MethodGet, "/tasks/status", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.GetTaskCountByStatus(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to retrieve task counts", response["error"])
	})
}

func (s *TaskHandlerSuite) resetMocks() {
	s.mockUserUsecase.ExpectedCalls = nil
	s.mockUserUsecase.Calls = nil
	s.mockTaskUsecase.ExpectedCalls = nil
	s.mockTaskUsecase.Calls = nil

}
