package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/yiheyistm/task_manager/internal/domain"
)

type TaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUseCase(taskRepo domain.TaskRepository) domain.ITaskUseCase {
	return &TaskUseCase{taskRepo: taskRepo}
}
func (uc *TaskUseCase) GetAll() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tasks, err := uc.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (uc *TaskUseCase) GetById(id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	task, err := uc.taskRepo.GetById(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}
func (uc *TaskUseCase) Create(task *domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if task == nil {
		return errors.New("task cannot be nil")
	}
	err := uc.taskRepo.Create(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (uc *TaskUseCase) Update(id string, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if task == nil {
		return errors.New("task cannot be nil")
	}
	err := uc.taskRepo.Update(ctx, id, task)
	if err != nil {
		return err
	}
	return nil
}

func (uc *TaskUseCase) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if id == "" {
		return errors.New("task ID cannot be empty")
	}
	err := uc.taskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *TaskUseCase) GetTasksByUser(username string) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	tasks, err := uc.taskRepo.GetByUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func (uc *TaskUseCase) GetTaskStatsByUser(username string) ([]domain.StatusCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	stats, err := uc.taskRepo.GetTaskStatsByUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
func (uc *TaskUseCase) GetTaskCountByStatus() ([]domain.StatusCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stats, err := uc.taskRepo.GetTaskCountByStatus(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (uc *TaskUseCase) GetByIdAndUser(id, username string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if id == "" || username == "" {
		return domain.Task{}, errors.New("task ID and username cannot be empty")
	}
	task, err := uc.taskRepo.GetByIdAndUser(ctx, id, username)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}
func (uc *TaskUseCase) UpdateByIdAndUser(id string, task *domain.Task, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if id == "" || username == "" {
		return errors.New("task ID and username cannot be empty")
	}
	if task == nil {
		return errors.New("task cannot be nil")
	}
	err := uc.taskRepo.UpdateByIdAndUser(ctx, id, task, username)
	if err != nil {
		return err
	}
	return nil
}
func (uc *TaskUseCase) DeleteByIdAndUser(id, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if id == "" || username == "" {
		return errors.New("task ID and username cannot be empty")
	}
	err := uc.taskRepo.DeleteByIdAndUser(ctx, id, username)
	if err != nil {
		return err
	}
	return nil
}
