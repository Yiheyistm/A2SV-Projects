package services

import (
	"context"
	"errors"
	"task_manager/config"
	"task_manager/models"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskServicesImpl struct {
	TaskCollection *mongo.Collection
}

func NewTaskServicesImpl() TaskServices {
	url := config.GetEnvString("MONGO_URI", "mongodb://localhost:27017")
	db := config.GetEnvString("MONGO_DB_NAME", "")
	collName := config.GetEnvString("MONGO_COLLECTION_NAME", "tasks")
	collection := config.ConnectDB(url, db, collName)
	config.CreateTaskIndexes(collection)
	return &TaskServicesImpl{
		TaskCollection: collection,
	}
}

func (s *TaskServicesImpl) GetAll() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []models.Task
	cursor, err := s.TaskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskServicesImpl) GetById(id string) (models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	filter := bson.M{"_id": objectID}
	if err := s.TaskCollection.FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

func (s *TaskServicesImpl) Create(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.TaskCollection.InsertOne(ctx, task)
	return err
}

func (s *TaskServicesImpl) Update(id string, updateTask *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": updateTask,
	}
	_, err = s.TaskCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (s *TaskServicesImpl) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.TaskCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (s *TaskServicesImpl) GetTaskCountByStatus() ([]models.StatusCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := s.TaskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.StatusCount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByUser
func (s *TaskServicesImpl) GetByUser(username string) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []models.Task
	filter := bson.M{"created_by": username}
	cursor, err := s.TaskCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByIdAndUser
func (s *TaskServicesImpl) GetByIdAndUser(taskID, username string) (models.Task, error) {
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	filter := bson.M{"_id": id, "created_by": username}
	if err := s.TaskCollection.FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

// UpdateByIdAndUser
func (s *TaskServicesImpl) UpdateByIdAndUser(id string, updateTask *models.Task, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": updateTask,
	}
	_, err = s.TaskCollection.UpdateOne(ctx, bson.M{"_id": objectID, "created_by": username}, update)
	return err
}

// DeleteByIdAndUser
func (s *TaskServicesImpl) DeleteByIdAndUser(id string, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.TaskCollection.DeleteOne(ctx, bson.M{"_id": objectID, "created_by": username})
	return err
}

// GetTaskStatsByUser
func (s *TaskServicesImpl) GetTaskStatsByUser(username string) ([]models.StatusCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_by": username}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := s.TaskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.StatusCount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
