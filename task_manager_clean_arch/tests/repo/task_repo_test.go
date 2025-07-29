package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/database"
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
	_ = godotenv.Load("../../../env") // for Testing purpose
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
	dbName := "test_db"
	s.database = s.client.Database(dbName)
	s.repository = persistence.NewTaskRepository(*s.database, "tasks")

	// Create a unique index on the id field
	indexModel := mongo.IndexModel{
		Keys: bson.M{"_id": 1}, // Index on the id field
	}
	_, err := s.database.Collection("tasks").Indexes().CreateOne(s.ctx, indexModel)
	if err != nil {
		s.T().Fatalf("Failed to create index: %v", err)
	}
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
	s.Run("Success", func() {
		collection := "tasks"
		repo := persistence.NewTaskRepository(*s.database, collection)

		taskRepo, ok := repo.(*persistence.TaskRepositoryImpl)
		s.True(ok)
		s.Equal(*s.database, taskRepo.Database)
		s.Equal(collection, taskRepo.Collection)
	})
}

// TestGetAll tests the GetAll method
func (s *TaskRepositorySuite) TestGetAllTasks() {
	s.Run("Success", func() {
		dueDate := time.Now().Add(48 * time.Hour).Truncate(time.Second).UTC()
		tasks := []database.TaskEntity{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Buy coffee from the store", CreatedBy: "Abebe", DueDate: primitive.NewDateTimeFromTime(dueDate), Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", Description: "Sell spices at the market", CreatedBy: "Kebede", DueDate: primitive.NewDateTimeFromTime(dueDate), Status: "completed"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetAll(s.ctx)
		taskList := database.FromTaskEntityListToDomainList(tasks)
		s.NoError(err)
		s.ElementsMatch(taskList, result)
	})

	s.Run("EmptyCollection", func() {
		_, err := s.database.Collection("tasks").DeleteMany(s.ctx, bson.M{})
		s.NoError(err)
		result, err := s.repository.GetAll(s.ctx)

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetById tests the GetById method
func (s *TaskRepositorySuite) TestGetById() {
	s.Run("Success", func() {
		// taskID := primitive.NewObjectID()
		dueDate := time.Now().Add(48 * time.Hour).Truncate(time.Second).UTC()
		taskEntity := database.TaskEntity{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Buy coffee from the store", CreatedBy: "Abebe", Status: "pending", DueDate: primitive.NewDateTimeFromTime(dueDate)}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, taskEntity)
		s.NoError(err)
		result, err := s.repository.GetById(s.ctx, taskEntity.ID.Hex())
		s.NoError(err)
		task := *database.FromTaskEntityToDomain(&taskEntity)
		task.CreatedBy = taskEntity.CreatedBy
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
	s.Run("Success", func() {
		dueDate := time.Now().Add(48 * time.Hour).Truncate(time.Second).UTC()
		taskEntity := database.TaskEntity{ID: primitive.NewObjectID(), Title: "Buy Coffee", Description: "Buy coffee from the store", CreatedBy: "Abebe", Status: "pending", DueDate: primitive.NewDateTimeFromTime(dueDate)}
		taskDomain := database.FromTaskEntityToDomain(&taskEntity)
		err := s.repository.Create(s.ctx, taskDomain)
		// Ensure the taskDomain is not nil
		s.NotNil(taskDomain)
		s.NoError(err)
		s.NotEqual(primitive.ObjectID{}, taskDomain.ID)

		// Verify the task was inserted
		var result database.TaskEntity
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskDomain.ID}).Decode(&result)
		s.NoError(err)
		s.Equal(taskEntity, result)
	})

	s.Run("DuplicateKey", func() {
		task := &database.TaskEntity{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		// Attempt to insert a task with the same ID
		duplicateTask := &database.TaskEntity{ID: task.ID, Title: "Sell Spices", CreatedBy: "Abebe", Status: "pending"}
		taskEntity := database.FromTaskEntityToDomain(duplicateTask)
		// Ensure the taskEntity is not nil
		s.NotNil(taskEntity)
		s.NoError(err)
		err = s.repository.Create(s.ctx, taskEntity)

		s.Error(err)
		s.Contains(err.Error(), "duplicate key")
	})
}

// TestUpdate tests the Update method
func (s *TaskRepositorySuite) TestUpdate() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		task := &database.TaskEntity{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		updatedTask := &database.TaskEntity{ID: taskID, Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		taskDomain := database.FromTaskEntityToDomain(updatedTask)
		err = s.repository.Update(s.ctx, taskID.Hex(), taskDomain)

		s.NoError(err)

		// Verify the update
		var result database.TaskEntity
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Decode(&result)
		s.NoError(err)
		s.Equal(updatedTask.Title, result.Title)
		s.Equal(updatedTask.Status, result.Status)
	})

	s.Run("InvalidID", func() {
		task := &database.TaskEntity{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		taskDomain := database.FromTaskEntityToDomain(task)
		err := s.repository.Update(s.ctx, "invalid_id", taskDomain)

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoUpdate", func() {
		taskID := primitive.NewObjectID()
		task := &database.TaskEntity{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		taskDomain := database.FromTaskEntityToDomain(task)
		err := s.repository.Update(s.ctx, taskID.Hex(), taskDomain)

		s.Error(err)
		s.Contains(err.Error(), "failed to update task or already up to date")
	})
}

// TestDelete tests the Delete method
func (s *TaskRepositorySuite) TestDelete() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		task := database.TaskEntity{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
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
	s.Run("Success", func() {
		tasks := []database.TaskEntity{
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
		_, err := s.database.Collection("tasks").DeleteMany(s.ctx, bson.M{})
		s.NoError(err)
		result, err := s.repository.GetTaskCountByStatus(s.ctx)

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetByUser tests the GetByUser method
func (s *TaskRepositorySuite) TestGetByUser() {
	s.Run("Success", func() {
		tasks := []database.TaskEntity{
			{ID: primitive.NewObjectID(), Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"},
			{ID: primitive.NewObjectID(), Title: "Sell Spices", CreatedBy: "Abebe", Status: "completed"},
			{ID: primitive.NewObjectID(), Title: "Brew Coffee", CreatedBy: "Kebede", Status: "pending"},
		}
		for _, task := range tasks {
			_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
			s.NoError(err)
		}

		result, err := s.repository.GetByUser(s.ctx, "Abebe")
		taskList := database.FromTaskEntityListToDomainList(tasks[:2]) // Only Abebe's tasks
		s.NoError(err)
		s.Len(result, 2)
		s.ElementsMatch(taskList, result)
	})

	s.Run("NoTasks", func() {
		result, err := s.repository.GetByUser(s.ctx, "Samson")

		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetByIdAndUser tests the GetByIdAndUser method
func (s *TaskRepositorySuite) TestGetByIdAndUser() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		task := database.TaskEntity{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		result, err := s.repository.GetByIdAndUser(s.ctx, taskID.Hex(), "Abebe")
		taskDomain := database.FromTaskEntityToDomain(&task)
		s.NoError(err)
		s.Equal(taskDomain, &result)
	})

	s.Run("InvalidID", func() {
		result, err := s.repository.GetByIdAndUser(s.ctx, "invalid_id", "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
		s.Equal(domain.Task{}, result)
	})

	s.Run("TaskNotFound", func() {
		taskID := primitive.NewObjectID()
		result, err := s.repository.GetByIdAndUser(s.ctx, taskID.Hex(), "samson")

		s.Error(err)
		s.Contains(err.Error(), "task not found")
		s.Equal(domain.Task{}, result)
	})
}

// TestUpdateByIdAndUser tests the UpdateByIdAndUser method
func (s *TaskRepositorySuite) TestUpdateByIdAndUser() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		task := &database.TaskEntity{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		updatedTask := &database.TaskEntity{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		updatedTaskDomain := database.FromTaskEntityToDomain(updatedTask)
		err = s.repository.UpdateByIdAndUser(s.ctx, taskID.Hex(), updatedTaskDomain, "Abebe")

		s.NoError(err)
		s.Equal(taskID, updatedTaskDomain.ID)

		// Verify the update
		var result database.TaskEntity
		err = s.database.Collection("tasks").FindOne(s.ctx, bson.M{"_id": taskID}).Decode(&result)
		s.NoError(err)
		s.Equal(updatedTaskDomain.Title, result.Title)
		s.Equal(updatedTaskDomain.Status, result.Status)
	})

	s.Run("InvalidID", func() {
		task := &database.TaskEntity{Title: "Buy Spices", CreatedBy: "Abebe", Status: "completed"}
		taskDomain := database.FromTaskEntityToDomain(task)
		err := s.repository.UpdateByIdAndUser(s.ctx, "invalid_id", taskDomain, "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "invalid ObjectID")
	})

	s.Run("NoUpdate", func() {
		taskID := primitive.NewObjectID()
		task := &database.TaskEntity{ID: taskID, Title: "Buy Spices", CreatedBy: "Abebe", Status: "pending"}
		_, err := s.database.Collection("tasks").InsertOne(s.ctx, task)
		s.NoError(err)

		updatedTask := &database.TaskEntity{Title: "Buy Spices", CreatedBy: "Abebe", Status: "pending"}
		taskDomain := database.FromTaskEntityToDomain(updatedTask)
		err = s.repository.UpdateByIdAndUser(s.ctx, taskID.Hex(), taskDomain, "Abebe")

		s.Error(err)
		s.Contains(err.Error(), "failed to update task or task not found for user")
	})
}

// TestDeleteByIdAndUser tests the DeleteByIdAndUser method
func (s *TaskRepositorySuite) TestDeleteByIdAndUser() {
	s.Run("Success", func() {
		taskID := primitive.NewObjectID()
		task := database.TaskEntity{ID: taskID, Title: "Buy Coffee", CreatedBy: "Abebe", Status: "pending"}
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
	s.Run("Success", func() {
		tasks := []database.TaskEntity{
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
		result, err := s.repository.GetTaskStatsByUser(s.ctx, "samson")

		s.NoError(err)
		s.Empty(result)
	})
}
