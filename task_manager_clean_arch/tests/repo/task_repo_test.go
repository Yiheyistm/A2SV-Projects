package persistence

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskRepositorySuite defines the test suite for taskRepository
type TaskRepositorySuite struct {
	suite.Suite
	client     *mongo.Client
	database   *mongo.Database
	repository domain.TaskRepository
	ctx        context.Context
}

// SetupSuite connects to MongoDB and initializes the client
func (s *TaskRepositorySuite) SetupSuite() {
	env := config.Load()
	DBHostURI := fmt.Sprintf("mongodb+srv://%s:%s@%s.r31b5bc.mongodb.net/?retryWrites=true&w=majority", env.DBUser, env.DBPass, env.DBHost)
	var err error
	s.ctx = context.Background()
	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(DBHostURI))
	if err != nil {
		s.T().Fatalf("Failed to connect to MongoDB: %v", err)
	}
}

// TearDownSuite disconnects the MongoDB client
func (s *TaskRepositorySuite) TearDownSuite() {
	if err := s.client.Disconnect(s.ctx); err != nil {
		s.T().Fatalf("Failed to disconnect MongoDB client: %v", err)
	}
}

// SetupTest initializes a unique test database and repository
func (s *TaskRepositorySuite) SetupTest() {
	// Create a unique database name using a timestamp
	dbName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	s.database = s.client.Database(dbName)
	s.repository = persistence.NewTaskRepository(*s.database, "tasks")
	s.ctx = context.Background()
}

// TearDownTest drops the test database to ensure isolation
func (s *TaskRepositorySuite) TearDownTest() {
	if err := s.database.Drop(s.ctx); err != nil {
		s.T().Fatalf("Failed to drop test database: %v", err)
	}
}

// TestTaskRepositorySuite runs the test suite
func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}

// TestNewTaskRepository tests the NewTaskRepository function
func (s *TaskRepositorySuite) TestNewTaskRepository() {
	s.Run("MerkatoSuccess", func() {
		collection := "tasks"
		repo := persistence.NewTaskRepository(*s.database, collection)

		taskRepo, ok := repo.(*persistence.TaskRepositoryImpl)
		s.True(ok)
		s.Equal(*s.database, taskRepo.Database)
		s.Equal(collection, taskRepo.Collection)
	})
}

// TestGetAll tests the GetAll method
func (s *TaskRepositorySuite) TestGetAll() {
	s.Run("MerkatoSuccess", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", CreatedBy: "Kebede", Status: "completed"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetAll(s.ctx)

		s.NoError(err)
		s.ElementsMatch(tasks, result)
	})

	s.Run("EmptyCollection", func() {
		result, err := s.repository.GetAll(s.ctx)

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetById tests the GetById method
func (s *TaskRepositorySuite) TestGetById() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		result, err := s.repository.GetById(s.ctx, taskID.Hex())

		s.NoError(err)
		s.Equal(task, result)
	})

	s.Run("InvalidID", func() {
		result, err := s.repository.GetById(s.ctx, "invalid_id")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
		s.Equal(domain.Task{}, result)
	})

	s.Run("TaskNotFound", func() {
		taskID := primitive.NewObjectID()
		result, err := s.repository.GetById(s.ctx, taskID.Hex())

		s.Error(err)
		s.Contains(err.Error(), "task not found")
		s.Equal(domain.Task{}, result)
	})
}

// TestCreate tests the Create method
func (s *TaskRepositorySuite) TestCreate() {
	s.Run("MerkatoSuccess", func() {
		task := &domain.Task{Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}

		err := s.repository.Create(s.ctx, task)

		s.NoError(err)
		s.NotEqual(primitive.ObjectID{}, task.ID)

		// Verify the task was inserted
		var result domain.Task
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": task.ID}).Decode(&result)
		s.NoError(err)
		s.Equal(task, &result)
	})

	s.Run("DuplicateKey", func() {
		task := &domain.Task{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		// Attempt to insert a task with the same ID
		duplicateTask := &domain.Task{ID: task.ID, Title: "Sell Spices", CreatedBy: "Abebe", Status: "pending"}
		err = s.repository.Create(s.ctx, duplicateTask)

		s.Error(err)
		s.Contains(err.Error(), "duplicate key")
	})
}

// TestUpdate tests the Update method
func (s *TaskRepositorySuite) TestUpdate() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := &domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		updatedTask := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err = s.repository.Update(s.ctx, taskID.Hex(), updatedTask)

		s.NoError(err)

		// Verify the update
		var result domain.Task
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Decode(&result)
		s.NoError(err)
		s.Equal(updatedTask.Title, result.Title)
		s.Equal(updatedTask.Status, result.Status)
	})

	s.Run("InvalidID", func() {
		task := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err := s.repository.Update(s.ctx, "invalid_id", task)

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoUpdate", func() {
		taskID := primitive.NewObjectID()
		task := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err := s.repository.Update(s.ctx, taskID.Hex(), task)

		s.Error(err)
		s.Contains(err.Error(), "failed to update task or already up to date")
	})
}

// TestDelete tests the Delete method
func (s *TaskRepositorySuite) TestDelete() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		err = s.repository.Delete(s.ctx, taskID.Hex())

		s.NoError(err)

		// Verify the task was deleted
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Err()
		s.Error(err)
		s.Equal(mongo.ErrNoDocuments, err)
	})

	s.Run("InvalidID", func() {
		err := s.repository.Delete(s.ctx, "invalid_id")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoDelete", func() {
		taskID := primitive.NewObjectID()
		err := s.repository.Delete(s.ctx, taskID.Hex())

		s.Error(err)
		s.Contains(err.Error(), "task not found or already deleted")
	})
}

// TestGetTaskCountByStatus tests the GetTaskCountByStatus method
func (s *TaskRepositorySuite) TestGetTaskCountByStatus() {
	s.Run("MerkatoSuccess", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Brew Coffee", CreatedBy: "Kebede", Status: "completed"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetTaskCountByStatus(s.ctx)

		s.NoError(err)
		expected := []domain.StatusCount{
			{Status: "pending", Count: 2},
			{Status: "completed", Count: 1},
		}
		s.ElementsMatch(expected, result)
	})

	s.Run("EmptyCollection", func() {
		result, err := s.repository.GetTaskCountByStatus(s.ctx)

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetByUser tests the GetByUser method
func (s *TaskRepositorySuite) TestGetByUser() {
	s.Run("MerkatoSuccess", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", CreatedBy: "Abebe", Status: "completed"},
			{ID: primitive.NewObjectID(), Title: "Brew Coffee", CreatedBy: "Kebede", Status: "pending"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetByUser(s.ctx, "Abebe")

		s.NoError(err)
		s.Len(result, 2)
		s.ElementsMatch([]domain.Task{tasks[0], tasks[1]}, result)
	})

	s.Run("NoTasks", func() {
		result, err := s.repository.GetByUser(s.ctx, "Abebe")

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetByIdAndUser tests the GetByIdAndUser method
func (s *TaskRepositorySuite) TestGetByIdAndUser() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		result, err := s.repository.GetByIdAndUser(s.ctx, taskID.Hex(), "Abebe")

		s.NoError(err)
		s.Equal(task, result)
	})

	s.Run("InvalidID", func() {
		result, err := s.repository.GetByIdAndUser(s.ctx, "invalid_id", "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
		s.Equal(domain.Task{}, result)
	})

	s.Run("TaskNotFound", func() {
		taskID := primitive.NewObjectID()
		result, err := s.repository.GetByIdAndUser(s.ctx, taskID.Hex(), "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "task not found")
		s.Equal(domain.Task{}, result)
	})
}

// TestUpdateByIdAndUser tests the UpdateByIdAndUser method
func (s *TaskRepositorySuite) TestUpdateByIdAndUser() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := &domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		updatedTask := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err = s.repository.UpdateByIdAndUser(s.ctx, taskID.Hex(), updatedTask, "Abebe")

		s.NoError(err)
		s.Equal(taskID, updatedTask.ID)

		// Verify the update
		var result domain.Task
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Decode(&result)
		s.NoError(err)
		s.Equal(updatedTask.Title, result.Title)
		s.Equal(updatedTask.Status, result.Status)
	})

	s.Run("InvalidID", func() {
		task := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err := s.repository.UpdateByIdAndUser(s.ctx, "invalid_id", task, "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoUpdate", func() {
		taskID := primitive.NewObjectID()
		task := &domain.Task{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		err := s.repository.UpdateByIdAndUser(s.ctx, taskID.Hex(), task, "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "failed to update task or task not found for user")
	})
}

// TestDeleteByIdAndUser tests the DeleteByIdAndUser method
func (s *TaskRepositorySuite) TestDeleteByIdAndUser() {
	s.Run("MerkatoSuccess", func() {
		taskID := primitive.NewObjectID()
		task := domain.Task{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		err = s.repository.DeleteByIdAndUser(s.ctx, taskID.Hex(), "Abebe")

		s.NoError(err)

		// Verify the task was deleted
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Err()
		s.Error(err)
		s.Equal(mongo.ErrNoDocuments, err)
	})

	s.Run("InvalidID", func() {
		err := s.repository.DeleteByIdAndUser(s.ctx, "invalid_id", "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoDelete", func() {
		taskID := primitive.NewObjectID()
		err := s.repository.DeleteByIdAndUser(s.ctx, taskID.Hex(), "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "task not found or not owned by user")
	})
}

// TestGetTaskStatsByUser tests the GetTaskStatsByUser method
func (s *TaskRepositorySuite) TestGetTaskStatsByUser() {
	s.Run("MerkatoSuccess", func() {
		tasks := []domain.Task{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Brew Coffee", CreatedBy: "Abebe", Status: "completed"},
			{ID: primitive.NewObjectID(), Title: "Deliver Goods", CreatedBy: "Kebede", Status: "pending"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetTaskStatsByUser(s.ctx, "Abebe")

		s.NoError(err)
		expected := []domain.StatusCount{
			{Status: "pending", Count: 2},
			{Status: "completed", Count: 1},
		}
		s.ElementsMatch(expected, result)
	})

	s.Run("NoTasks", func() {
		result, err := s.repository.GetTaskStatsByUser(s.ctx, "Abebe")

		s.NoError(err)
		s.Empty(result)
	})
}
