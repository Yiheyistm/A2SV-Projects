package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskMapperSuite defines the test suite for task_mapper.go
type TaskMapperSuite struct {
	suite.Suite
}

// TestTaskMapperSuite runs the test suite
func TestTaskMapperSuite(t *testing.T) {
	suite.Run(t, new(TaskMapperSuite))
}

// TestFromRequestToDomainTask tests the FromRequestToDomainTask function
func (s *TaskMapperSuite) TestFromRequestToDomainTask() {
	s.Run("Success", func() {
		taskRequest := dto.TaskRequest{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedBy:   "Abebe",
		}
		expectedTask := &domain.Task{
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedBy:   "Abebe",
		}

		result := taskRequest.FromRequestToDomainTask()

		s.Equal(expectedTask, result)
	})

	s.Run("EmptyFields", func() {
		taskRequest := dto.TaskRequest{
			Title:       "",
			Description: "",
			DueDate:     time.Time{},
			Status:      "",
			CreatedBy:   "",
		}
		expectedTask := &domain.Task{
			Title:       "",
			Description: "",
			DueDate:     time.Time{},
			Status:      "",
			CreatedBy:   "",
		}

		result := taskRequest.FromRequestToDomainTask()

		s.Equal(expectedTask, result)
	})
}

// TestFromDomainTaskToResponse tests the dto.FromDomainTaskToResponse function
func (s *TaskMapperSuite) TestdFromDomainTaskToResponse() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		domainTask := &domain.Task{
			ID:          taskID,
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedBy:   "Abebe",
		}
		expectedResponse := &dto.TaskResponse{
			ID:          taskID.Hex(),
			Title:       "Buy Coffee",
			Description: "Get buna from Merkato",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      "pending",
			CreatedBy:   "Abebe",
		}

		result := dto.FromDomainTaskToResponse(domainTask)

		s.Equal(expectedResponse, result)
	})

	s.Run("EmptyFields", func() {
		domainTask := &domain.Task{
			ID:          primitive.ObjectID{},
			Title:       "",
			Description: "",
			DueDate:     time.Time{},
			Status:      "",
			CreatedBy:   "",
		}
		expectedResponse := &dto.TaskResponse{
			ID:          primitive.ObjectID{}.Hex(),
			Title:       "",
			Description: "",
			DueDate:     time.Time{},
			Status:      "",
			CreatedBy:   "",
		}

		result := dto.FromDomainTaskToResponse(domainTask)

		s.Equal(expectedResponse, result)
	})

	s.Run("NilInput", func() {
		var domainTask domain.Task
		expectedResponse := &dto.TaskResponse{
			ID: "000000000000000000000000",
		}

		result := dto.FromDomainTaskToResponse(&domainTask)

		s.Equal(expectedResponse, result)
	})
}

// Testdto.FromDomainTaskToResponseList tests the dto.FromDomainTaskToResponseList function
func (s *TaskMapperSuite) TestFromDomainTaskToResponseList() {
	s.Run("Success", func() {
		taskID1 := primitive.NewObjectID()
		taskID2 := primitive.NewObjectID()
		domainTasks := []domain.Task{
			{
				ID:          taskID1,
				Title:       "Buy Coffee",
				Description: "Get buna from Merkato",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "pending",
				CreatedBy:   "Abebe",
			},
			{
				ID:          taskID2,
				Title:       "Sell Spices",
				Description: "Trade in Merkato",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "completed",
				CreatedBy:   "Kebede",
			},
		}
		expectedResponses := []dto.TaskResponse{
			{
				ID:          taskID1.Hex(),
				Title:       "Buy Coffee",
				Description: "Get buna from Merkato",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "pending",
				CreatedBy:   "Abebe",
			},
			{
				ID:          taskID2.Hex(),
				Title:       "Sell Spices",
				Description: "Trade in Merkato",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "completed",
				CreatedBy:   "Kebede",
			},
		}

		result := dto.FromDomainTaskToResponseList(domainTasks)

		s.Equal(expectedResponses, result)
	})

	s.Run("EmptySlice", func() {
		domainTasks := []domain.Task{}
		var expectedResponses []dto.TaskResponse

		result := dto.FromDomainTaskToResponseList(domainTasks)

		s.Equal(expectedResponses, result)
	})

	s.Run("NilSlice", func() {
		var domainTasks []domain.Task
		var expectedResponses []dto.TaskResponse

		result := dto.FromDomainTaskToResponseList(domainTasks)

		s.Equal(expectedResponses, result)
	})
}
