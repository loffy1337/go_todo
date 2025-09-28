package repository

import (
	"context"

	"github.com/loffy1337/go_todo/internal/models"
)

// Описание операций с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateProfile(ctx context.Context, id uint64, avatarURL *string, statusText string) error
	UpdatePassword(ctx context.Context, id uint64, newPasswordHash string) error
	SetOnline(ctx context.Context, id uint64, online bool) error
	UpdateLastActive(ctx context.Context, id uint64) error
}

// Описание операций с группами задач
type TaskGroupRepository interface {
	Create(ctx context.Context, taskGroup *models.TaskGroup) (uint64, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	UpdateTitle(ctx context.Context, id uint64, userID uint64, newTitle string) error
	ListByUser(ctx context.Context, userID uint64) ([]models.TaskGroup, error)
}

// Описание операций с задачами
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (uint64, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	Update(ctx context.Context, id uint64, userID uint64, title string, body string) error
	MoveToGroup(ctx context.Context, id uint64, userID uint64, newGroupID uint64) error
	SetDone(ctx context.Context, id uint64, userID uint64, done bool) error
	ListByGroup(ctx context.Context, userID uint64, groupID uint64) ([]models.Task, error)
}
