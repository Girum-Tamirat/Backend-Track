package Usecases

import (
	"context"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecase interface {
	CreateTask(t Domain.Task) (Domain.Task, error)
	ListTasks() ([]Domain.Task, error)
	GetTaskByID(id primitive.ObjectID) (Domain.Task, error)
	UpdateTask(id primitive.ObjectID, patch map[string]interface{}) (Domain.Task, error)
	DeleteTask(id primitive.ObjectID) error
}

type taskUsecase struct {
	repo Repositories.TaskRepository
	timeout time.Duration
}

func NewTaskUsecase(r Repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{repo: r, timeout: 5 * time.Second}
}

func (u *taskUsecase) CreateTask(t Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.Create(ctx, t)
}

func (u *taskUsecase) ListTasks() ([]Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.FindAll(ctx)
}

func (u *taskUsecase) GetTaskByID(id primitive.ObjectID) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.FindByID(ctx, id)
}

func (u *taskUsecase) UpdateTask(id primitive.ObjectID, patch map[string]interface{}) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.Update(ctx, id, patch)
}

func (u *taskUsecase) DeleteTask(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.Delete(ctx, id)
}
