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
	collName := config.GetEnvString("MONGO_COLLECTION_NAME", "users")
	collection := config.ConnectDB(url, db, collName)
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
