package persistence

import (
	"context"
	"errors"

	"github.com/yiheyistm/task_manager/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepositoryImpl struct {
	Database   mongo.Database
	Collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &TaskRepositoryImpl{
		Database:   db,
		Collection: collection,
	}
}

func (r *TaskRepositoryImpl) GetAll(ctx context.Context) ([]domain.Task, error) {

	var tasks []domain.Task
	cursor, err := r.Database.Collection(r.Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (r *TaskRepositoryImpl) GetById(ctx context.Context, id string) (domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, err
	}

	var task domain.Task
	filter := bson.M{"_id": objectID}
	if err := r.Database.Collection(r.Collection).FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("task not found")
		}
		return domain.Task{}, err
	}
	return task, nil
}

func (s *TaskRepositoryImpl) Create(ctx context.Context, task *domain.Task) error {
	result, err := s.Database.Collection(s.Collection).InsertOne(ctx, task)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		return errors.New("failed to insert task")
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *TaskRepositoryImpl) Update(ctx context.Context, id string, updateTask *domain.Task) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": updateTask,
	}
	result, err := s.Database.Collection(s.Collection).UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("failed to update task or already up to date")
	}
	return nil
}

func (s *TaskRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := s.Database.Collection(s.Collection).DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found or already deleted")
	}
	return nil
}

func (s *TaskRepositoryImpl) GetTaskCountByStatus(ctx context.Context) ([]domain.StatusCount, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := s.Database.Collection(s.Collection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []domain.StatusCount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByUser
func (s *TaskRepositoryImpl) GetByUser(ctx context.Context, username string) ([]domain.Task, error) {

	var tasks []domain.Task
	filter := bson.M{"created_by": username}
	cursor, err := s.Database.Collection(s.Collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByIdAndUser
func (s *TaskRepositoryImpl) GetByIdAndUser(ctx context.Context, taskID, username string) (domain.Task, error) {
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	var task domain.Task
	filter := bson.M{"_id": id, "created_by": username}
	if err := s.Database.Collection(s.Collection).FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("task not found")
		}
		return domain.Task{}, err
	}
	return task, nil
}

// UpdateByIdAndUser
func (s *TaskRepositoryImpl) UpdateByIdAndUser(ctx context.Context, id string, updateTask *domain.Task, username string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": updateTask,
	}
	result, err := s.Database.Collection(s.Collection).UpdateOne(ctx, bson.M{"_id": objectID, "created_by": username}, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("failed to update task or task not found for user")
	}
	updateTask.ID = objectID
	return nil
}

// DeleteByIdAndUser
func (s *TaskRepositoryImpl) DeleteByIdAndUser(ctx context.Context, id string, username string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := s.Database.Collection(s.Collection).DeleteOne(ctx, bson.M{"_id": objectID, "created_by": username})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found or not owned by user")
	}
	return nil
}

// GetTaskStatsByUser
func (s *TaskRepositoryImpl) GetTaskStatsByUser(ctx context.Context, username string) ([]domain.StatusCount, error) {

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_by": username}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := s.Database.Collection(s.Collection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []domain.StatusCount
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
