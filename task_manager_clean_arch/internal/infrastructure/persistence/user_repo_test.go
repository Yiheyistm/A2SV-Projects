package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/config"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/infrastructure/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepositorySuite defines the test suite for userRepository
type UserRepositorySuite struct {
	suite.Suite
	client     *mongo.Client
	database   *mongo.Database
	repository domain.UserRepository
	ctx        context.Context
}

// SetupSuite connects to MongoDB and initializes the client
func (s *UserRepositorySuite) SetupSuite() {
	godotenv.Load("../../env")
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
func (s *UserRepositorySuite) TearDownSuite() {
	if err := s.client.Disconnect(s.ctx); err != nil {
		s.T().Fatalf("Failed to disconnect MongoDB client: %v", err)
	}
}

// SetupTest initializes a unique test database and repository
func (s *UserRepositorySuite) SetupTest() {
	godotenv.Load("../../env")
	dbName := "test_db"
	s.database = s.client.Database(dbName)
	s.repository = NewUserRepository(*s.database, "users")

	// Create a unique index on the username field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Unique on the username field
		Options: options.Index().SetUnique(true),
	}
	_, err := s.database.Collection("users").Indexes().CreateOne(s.ctx, indexModel)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		s.T().Fatalf("Failed to create unique index: %v", err)
	}
	if err != nil {
		s.T().Fatalf("Failed to create unique index: %v", err)
	}
	s.ctx = context.Background()
}

// TearDownTest drops the test database to ensure isolation
func (s *UserRepositorySuite) TearDownTest() {
	if err := s.database.Drop(s.ctx); err != nil {
		s.T().Fatalf("Failed to drop test database: %v", err)
	}
}

// TestUserRepositorySuite runs the test suite
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

// TestNewUserRepository tests the NewUserRepository function
func (s *UserRepositorySuite) TestNewUserRepository() {
	s.Run("Success", func() {
		collection := "users"
		repo := NewUserRepository(*s.database, collection)

		userRepo, ok := repo.(*UserRepositoryImpl)
		s.True(ok)
		s.Equal(*s.database, userRepo.DB)
		s.Equal(collection, userRepo.Collection)
	})
}

// TestInsert tests the Insert method
func (s *UserRepositorySuite) TestUserInsert() {
	s.Run("Success", func() {
		id := primitive.NewObjectID()
		user := &database.UserEntity{
			ID:       id,
			Username: "Abebe" + id.Hex(), // to remove duplicate key error
			Email:    "abebe@example.com",
			Password: "hashed_password",
			Role:     "user",
		}

		// Use the ObjectID directly for the _id field in MongoDB
		res, err := s.database.Collection("users").InsertOne(s.ctx, user)
		s.NoError(err)
		s.NotNil(res.InsertedID)
		s.NoError(err)

		// Verify the user was inserted
		var result database.UserEntity
		err = s.database.Collection("users").FindOne(s.ctx, bson.M{"username": "Abebe" + id.Hex()}).Decode(&result)
		s.NoError(err)
		s.Equal(user, &result)
	})

	s.Run("NilUser", func() {
		err := s.repository.Insert(s.ctx, nil)
		s.Error(err)
		s.Contains(err.Error(), "user cannot be nil")
	})

	s.Run("DuplicateKey", func() {
		user := &domain.User{
			ID:       "1",
			Username: "Nicko",
			Email:    "nicko@example.com",
			Password: "hashed_password",
			Role:     "user",
		}
		err := s.repository.Insert(s.ctx, user)
		s.NoError(err)

		// Insert again to trigger duplicate key error
		err = s.repository.Insert(s.ctx, user)

		s.Error(err)
		s.Contains(err.Error(), "duplicate key")
	})
}

// TestGetAll tests the GetAll method
func (s *UserRepositorySuite) TestGetAll() {
	s.Run("Success", func() {
		users := []database.UserEntity{
			{ID: primitive.NewObjectID(), Username: "Abebe", Email: "abebe@example.com", Role: "user"},
			{ID: primitive.NewObjectID(), Username: "Kebede", Email: "kebede@example.com", Role: "admin"},
		}
		for _, user := range users {
			_, err := s.database.Collection("users").InsertOne(s.ctx, user)
			s.NoError(err)
		}

		result, err := s.repository.GetAll(s.ctx)
		userDomain := database.FromEntityListToDomainList(users)
		s.NoError(err)
		s.ElementsMatch(userDomain, result)
	})

	s.Run("EmptyCollection", func() {
		_, err := s.database.Collection("users").DeleteMany(s.ctx, bson.M{})
		s.NoError(err)
		result, err := s.repository.GetAll(s.ctx)
		s.NoError(err)
		s.Empty(result)
	})
}

// TestGetUser tests the getUser helper method
func (s *UserRepositorySuite) TestGetUser() {
	s.Run("Success", func() {
		user := &database.UserEntity{
			ID:       primitive.NewObjectID(),
			Username: "Abebe",
			Email:    "abebe@example.com",
			Role:     "user",
		}
		_, err := s.database.Collection("users").InsertOne(s.ctx, user)
		s.NoError(err)

		result, err := s.repository.GetUser(s.ctx, "username", "Abebe")
		userDomain := database.FromEntityToDomain(user)
		s.NoError(err)
		s.Equal(userDomain, result)
	})

	s.Run("UserNotFound", func() {
		result, err := s.repository.GetUser(s.ctx, "username", "Kebede")
		s.Error(err)
		s.Contains(err.Error(), "user not found")
		s.Nil(result)
	})
}

// TestGetByUsername tests the GetByUsername method
func (s *UserRepositorySuite) TestGetByUsername() {
	s.Run("Success", func() {
		user := &database.UserEntity{
			ID:       primitive.NewObjectID(),
			Username: "Abebe",
			Email:    "abebe@example.com",
			Role:     "user",
		}
		_, err := s.database.Collection("users").InsertOne(s.ctx, user)
		s.NoError(err)

		result, err := s.repository.GetByUsername(s.ctx, "Abebe")
		userDomain := database.FromEntityToDomain(user)
		s.NoError(err)
		s.Equal(userDomain, result)
	})

	s.Run("UserNotFound", func() {
		result, err := s.repository.GetByUsername(s.ctx, "Kebede")

		s.Error(err)
		s.Contains(err.Error(), "user not found")
		s.Nil(result)
	})
}

// TestGetByEmail tests the GetByEmail method
func (s *UserRepositorySuite) TestGetByEmail() {
	s.Run("Success", func() {
		user := &database.UserEntity{
			ID:       primitive.NewObjectID(),
			Username: "Abebe",
			Email:    "abebe@example.com",
			Role:     "user",
		}
		_, err := s.database.Collection("users").InsertOne(s.ctx, user)
		s.NoError(err)
		result, err := s.repository.GetByEmail(s.ctx, "abebe@example.com")
		s.NoError(err)
		userDomain := database.FromEntityToDomain(user)
		s.Equal(userDomain, result)
	})

	s.Run("UserNotFound", func() {
		result, err := s.repository.GetByEmail(s.ctx, "kebede@example.com")

		s.Error(err)
		s.Contains(err.Error(), "user not found")
		s.Nil(result)
	})
}

// TestGetUserFromContext tests the GetUserFromContext method
func (s *UserRepositorySuite) TestGetUserFromContext() {
	s.Run("Success", func() {
		user := &database.UserEntity{
			ID:       primitive.NewObjectID(),
			Username: "Abebe",
			Email:    "abebe@example.com",
			Role:     "user",
		}
		_, err := s.database.Collection("users").InsertOne(s.ctx, user)
		s.NoError(err)

		c, _ := gin.CreateTestContext(nil)
		c.Set("username", "Abebe")

		result := s.repository.GetUserFromContext(c)
		userDomain := database.FromEntityToDomain(user)
		s.Equal(userDomain, result)
	})

	s.Run("UserNotFound", func() {
		c, _ := gin.CreateTestContext(nil)
		c.Set("username", "Kebede")

		result := s.repository.GetUserFromContext(c)

		s.Nil(result)
	})

	s.Run("EmptyUsername", func() {
		c, _ := gin.CreateTestContext(nil)
		// No username set in context

		result := s.repository.GetUserFromContext(c)

		s.Nil(result)
	})
}
