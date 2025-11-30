package Repositories

import (
	"context"

	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository defines task persistence operations
type TaskRepository interface {
	Create(ctx context.Context, t Domain.Task) (Domain.Task, error)
	FindAll(ctx context.Context) ([]Domain.Task, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (Domain.Task, error)
	Update(ctx context.Context, id primitive.ObjectID, patch map[string]interface{}) (Domain.Task, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// UserRepository defines user persistence operations
type UserRepository interface {
	Create(ctx context.Context, u Domain.User) (Domain.User, error)
	FindByUsername(ctx context.Context, username string) (Domain.User, error)
	UpdateRole(ctx context.Context, username, role string) error
	CountUsers(ctx context.Context) (int64, error)
}
