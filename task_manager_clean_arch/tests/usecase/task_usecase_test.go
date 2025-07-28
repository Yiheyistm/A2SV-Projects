package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/usecase"
	mocks_domain "github.com/yiheyistm/task_manager/mocks/mocks_domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskUseCaseSuite defines the test suite for TaskUseCase
type TaskUseCaseSuite struct {
	suite.Suite
	mockRepo *mocks_domain.TaskRepository
	useCase  domain.ITaskUseCase
}

// SetupTest initializes the mocks and use case before each test
func (s *TaskUseCaseSuite) SetupTest() {
	s.mockRepo = mocks_domain.NewTaskRepository(s.T())
	s.useCase = usecase.NewTaskUseCase(s.mockRepo)
}

// TestTaskUseCaseSuite runs the test suite
func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}

// TestGetAll tests the GetAll method
func (s *TaskUseCaseSuite) TestGetAll() {
	s.Run("Success", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", Description: "Trade in Merkato", Status: "completed", CreatedBy: "kebede", DueDate: time.Now()},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", Description: "Trade in Merkato", Status: "completed", CreatedBy: "kebede", DueDate: time.Now()},
		}
		s.mockRepo.On("GetAll", mock.Anything).Return(tasks, nil)

		result, err := s.useCase.GetAll()

		s.NoError(err)
		s.Equal(tasks, result)
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

// TestGetById tests the GetById method
func (s *TaskUseCaseSuite) TestGetById() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		task := domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("GetById", mock.Anything, id.Hex()).Return(task, nil)

		result, err := s.useCase.GetById(id.Hex())

		s.NoError(err)
		s.Equal(task, result)
	})

	s.Run("TaskNotFound", func() {
		s.mockRepo.On("GetById", mock.Anything, "unknown").Return(domain.Task{}, errors.New("task not found"))

		result, err := s.useCase.GetById("unknown")

		s.Error(err)
		s.EqualError(err, "task not found")
		s.Equal(domain.Task{}, result)
	})
}

// TestCreate tests the Create method
func (s *TaskUseCaseSuite) TestCreate() {
	s.Run("Success", func() {
		task := &domain.Task{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("Create", mock.Anything, task).Return(nil)
		err := s.useCase.Create(task)
		s.NoError(err)
	})
	s.Run("NilTask", func() {
		err := s.useCase.Create(nil)
		s.Error(err)
		s.EqualError(err, "task cannot be nil")
	})

	s.Run("RepositoryError", func() {
		task := &domain.Task{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("Create", mock.Anything, task).Return(errors.New("create failed"))
		err := s.useCase.Create(task)
		s.Error(err)
		s.EqualError(err, "create failed")
	})
}

// TestUpdate tests the Update method
func (s *TaskUseCaseSuite) TestUpdate() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		task := &domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("Update", mock.Anything, id.Hex(), task).Return(nil)
		err := s.useCase.Update(id.Hex(), task)
		s.NoError(err)
	})

	s.Run("NilTask", func() {
		id := primitive.NewObjectID()
		err := s.useCase.Update(id.Hex(), nil)

		s.Error(err)
		s.EqualError(err, "task cannot be nil")
	})

	s.Run("RepositoryError", func() {
		task := &domain.Task{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("Update", mock.Anything, task.ID.Hex(), task).Return(errors.New("update failed"))

		err := s.useCase.Update(task.ID.Hex(), task)

		s.Error(err)
		s.EqualError(err, "update failed")
	})
}

// TestDelete tests the Delete method
func (s *TaskUseCaseSuite) TestDelete() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		s.mockRepo.On("Delete", mock.Anything, id.Hex()).Return(nil)
		err := s.useCase.Delete(id.Hex())
		s.NoError(err)
	})

	s.Run("EmptyID", func() {
		err := s.useCase.Delete("")
		s.Error(err)
		s.EqualError(err, "task ID cannot be empty")
	})
	s.Run("RepositoryError", func() {
		id := primitive.NewObjectID()
		s.mockRepo.On("Delete", mock.Anything, id.Hex()).Return(errors.New("delete failed"))
		err := s.useCase.Delete(id.Hex())
		s.Error(err)
		s.EqualError(err, "delete failed")
	})
}

// TestGetTasksByUser tests the GetTasksByUser method
func (s *TaskUseCaseSuite) TestGetTasksByUser() {
	s.Run("Success", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", Description: "Trade in Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()},
		}
		s.mockRepo.On("GetByUser", mock.Anything, "abebe").Return(tasks, nil)

		result, err := s.useCase.GetTasksByUser("abebe")

		s.NoError(err)
		s.Equal(tasks, result)
	})

	s.Run("EmptyUsername", func() {
		result, err := s.useCase.GetTasksByUser("")
		s.Error(err)
		s.EqualError(err, "username cannot be empty")
		s.Nil(result)
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		s.mockRepo.On("GetByUser", mock.Anything, "abebe").Return(nil, errors.New("database error"))

		result, err := s.useCase.GetTasksByUser("abebe")
		s.Error(err)
		s.EqualError(err, "database error")
		s.Nil(result)
	})
}

// TestGetTaskStatsByUser tests the GetTaskStatsByUser method
func (s *TaskUseCaseSuite) TestGetTaskStatsByUser() {
	s.Run("Success", func() {
		stats := []domain.StatusCount{
			{Status: "pending", Count: 5},
			{Status: "completed", Count: 3},
		}
		s.mockRepo.On("GetTaskStatsByUser", mock.Anything, "abebe").Return(stats, nil)
		result, err := s.useCase.GetTaskStatsByUser("abebe")
		s.NoError(err)
		s.Equal(stats, result)
	})

	s.Run("EmptyUsername", func() {
		result, err := s.useCase.GetTaskStatsByUser("")
		s.Error(err)
		s.EqualError(err, "username cannot be empty")
		s.Nil(result)
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		s.mockRepo.On("GetTaskStatsByUser", mock.Anything, "abebe").Return(nil, errors.New("stats error"))
		result, err := s.useCase.GetTaskStatsByUser("abebe")
		s.Error(err)
		s.EqualError(err, "stats error")
		s.Nil(result)
	})
}

// TestGetTaskCountByStatus tests the GetTaskCountByStatus method
func (s *TaskUseCaseSuite) TestGetTaskCountByStatus() {
	s.Run("Success", func() {
		stats := []domain.StatusCount{
			{Status: "pending", Count: 10},
			{Status: "completed", Count: 8},
		}
		s.mockRepo.On("GetTaskCountByStatus", mock.Anything).Return(stats, nil)
		result, err := s.useCase.GetTaskCountByStatus()
		s.NoError(err)
		s.Equal(stats, result)
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		s.mockRepo.On("GetTaskCountByStatus", mock.Anything).Return(nil, errors.New("stats error"))
		result, err := s.useCase.GetTaskCountByStatus()
		s.Error(err)
		s.EqualError(err, "stats error")
		s.Nil(result)
	})
}

// TestGetByIdAndUser tests the GetByIdAndUser method
func (s *TaskUseCaseSuite) TestGetByIdAndUser() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		task := domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "pending", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("GetByIdAndUser", mock.Anything, task.ID.Hex(), "abebe").Return(task, nil)
		result, err := s.useCase.GetByIdAndUser(id.Hex(), "abebe")
		s.NoError(err)
		s.Equal(task, result)
	})

	s.Run("EmptyIDs", func() {
		result, err := s.useCase.GetByIdAndUser("", "")
		s.Error(err)
		s.EqualError(err, "task ID and username cannot be empty")
		s.Equal(domain.Task{}, result)
	})

	s.Run("RepositoryError", func() {
		id := primitive.NewObjectID()
		s.mockRepo.On("GetByIdAndUser", mock.Anything, id.Hex(), "abebe").Return(domain.Task{}, errors.New("task not found"))
		result, err := s.useCase.GetByIdAndUser(id.Hex(), "abebe")
		s.Error(err)
		s.EqualError(err, "task not found")
		s.Equal(domain.Task{}, result)
	})
}

// TestUpdateByIdAndUser tests the UpdateByIdAndUser method
func (s *TaskUseCaseSuite) TestUpdateByIdAndUser() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		task := &domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("UpdateByIdAndUser", mock.Anything, task.ID.Hex(), task, "abebe").Return(nil)
		err := s.useCase.UpdateByIdAndUser(id.Hex(), task, "abebe")
		s.NoError(err)
	})

	s.Run("EmptyIDs", func() {
		id := primitive.NewObjectID()
		task := &domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()}
		err := s.useCase.UpdateByIdAndUser("", task, "")
		s.Error(err)
		s.EqualError(err, "task ID and username cannot be empty")
	})

	s.Run("NilTask", func() {
		id := primitive.NewObjectID()
		err := s.useCase.UpdateByIdAndUser(id.Hex(), nil, "abebe")
		s.Error(err)
		s.EqualError(err, "task cannot be nil")
	})

	s.Run("RepositoryError", func() {
		s.mockRepo.ExpectedCalls = nil
		id := primitive.NewObjectID()
		task := &domain.Task{ID: id, Title: "Buy Coffee", Description: "Get buna from Merkato", Status: "completed", CreatedBy: "abebe", DueDate: time.Now()}
		s.mockRepo.On("UpdateByIdAndUser", mock.Anything, id.Hex(), task, "abebe").Return(errors.New("update failed"))
		err := s.useCase.UpdateByIdAndUser(id.Hex(), task, "abebe")
		s.Error(err)
		s.EqualError(err, "update failed")
	})
}

// TestDeleteByIdAndUser tests the DeleteByIdAndUser method
func (s *TaskUseCaseSuite) TestDeleteByIdAndUser() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		s.mockRepo.On("DeleteByIdAndUser", mock.Anything, id.Hex(), "abebe").Return(nil)
		err := s.useCase.DeleteByIdAndUser(id.Hex(), "abebe")

		s.NoError(err)
	})

	s.Run("EmptyIDs", func() {
		err := s.useCase.DeleteByIdAndUser("", "")
		s.EqualError(err, "task ID and username cannot be empty")
	})

	s.Run("RepositoryError", func() {
		id := primitive.NewObjectID()
		s.mockRepo.On("DeleteByIdAndUser", mock.Anything, id.Hex(), "abebe").Return(errors.New("delete failed"))
		err := s.useCase.DeleteByIdAndUser(id.Hex(), "abebe")
		s.Error(err)
		s.EqualError(err, "delete failed")
	})
}
