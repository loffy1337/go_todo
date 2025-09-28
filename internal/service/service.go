package service

import (
	"context"
	"fmt"
	"time"

	"github.com/loffy1337/go_todo/internal/models"
	"github.com/loffy1337/go_todo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrConflict     = fmt.Errorf("conflict")
	ErrNotFound     = fmt.Errorf("not found")
)

// Описание бизнес-логики авторизации
type AuthService interface {
	Register(ctx context.Context, username string, password string) (*models.User, error)
	Login(ctx context.Context, username string, password string) (*models.User, error)
	Logout(ctx context.Context, userID uint64) error
}

// Описание бизнес-логики профиля
type ProfileService interface {
	View(ctx context.Context, userID uint64) (*models.User, error)
	EditAvatarAndStatus(ctx context.Context, userID uint64, avatarURL *string, statusText string) error
	ChangePassword(ctx context.Context, userID uint64, oldPassword string, newPassword string) error
}

// Описание бизнес-логики группы задач
type TaskGroupService interface {
	Create(ctx context.Context, userID uint64, title string, status string) (*models.TaskGroup, error)
	Delete(ctx context.Context, userID uint64, groupID uint64) error
	Rename(ctx context.Context, userID uint64, groupID uint64, newTitle string) error
	ListByUser(ctx context.Context, userID uint64) ([]models.TaskGroup, error)
}

// Описание бизнес-логики задачи
type TaskService interface {
	Create(ctx context.Context, userID uint64, groupID uint64, title string, body string) (*models.Task, error)
	Edit(ctx context.Context, userID uint64, taskID uint64, title string, body string) error
	Delete(ctx context.Context, userID uint64, taskID uint64) error
	Move(ctx context.Context, userID uint64, taskID uint64, newGroupID uint64) error
	ToggleDone(ctx context.Context, userID uint64, taskID uint64, done bool) error
	ListByGroup(ctx context.Context, userID uint64, groupID uint64) ([]models.Task, error)
}

// Реализация каркаса сервисов

type Services struct {
	Auth       AuthService
	Profile    ProfileService
	TaskGroups TaskGroupService
	Tasks      TaskService
}

type store struct {
	users      repository.UserRepository
	taskGroups repository.TaskGroupRepository
	tasks      repository.TaskRepository
}

func New(s store) *Services {
	return &Services{
		Auth:       &authService{users: s.users},
		Profile:    &profileService{users: s.users},
		TaskGroups: &taskGroupService{taskGroup: s.taskGroups},
		Tasks:      &taskService{task: s.tasks},
	}
}

// Реализация бизнес-логики авторизации

type authService struct {
	users repository.UserRepository
}

func (service *authService) Register(ctx context.Context, username string, password string) (*models.User, error) {
	if len(username) < 3 || len(password) < 6 {
		return nil, fmt.Errorf("validation failed")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var user *models.User = &models.User{
		Username:     username,
		Password:     string(passwordHash),
		StatusText:   "",
		RegisteredAt: time.Now(),
		IsOnline:     false,
	}
	id, err := service.users.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (service *authService) Login(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := service.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrUnauthorized
	}
	err = service.users.SetOnline(ctx, user.ID, true)
	if err != nil {
		return nil, err
	}
	err = service.users.UpdateLastActive(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *authService) Logout(ctx context.Context, userID uint64) error {
	var err error = service.users.SetOnline(ctx, userID, false)
	if err != nil {
		return err
	}
	return nil
}

// Реализация бизнес-логики профиля

type profileService struct {
	users repository.UserRepository
}

func (service *profileService) View(ctx context.Context, userID uint64) (*models.User, error) {
	user, err := service.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *profileService) EditAvatarAndStatus(ctx context.Context, userID uint64, avatarURL *string, statusText string) error {
	var err error = service.users.UpdateProfile(ctx, userID, avatarURL, statusText)
	if err != nil {
		return err
	}
	return nil
}

func (service *profileService) ChangePassword(ctx context.Context, userID uint64, oldPassword string, newPassword string) error {
	user, err := service.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return ErrUnauthorized
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = service.users.UpdatePassword(ctx, userID, string(passwordHash))
	if err != nil {
		return err
	}
	return nil
}

// Реализация бизнес-логики группы задач

type taskGroupService struct {
	taskGroup repository.TaskGroupRepository
}

func (service *taskGroupService) Create(ctx context.Context, userID uint64, title string, status string) (*models.TaskGroup, error) {
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	if status != models.GroupStatusUrgent && status != models.GroupStatusDaily && status != models.GroupStatusLongTerm {
		return nil, fmt.Errorf("invalid group status")
	}
	var taskGroup *models.TaskGroup = &models.TaskGroup{
		UserID:    userID,
		Title:     title,
		Status:    status,
		CreatedAt: time.Now(),
	}
	id, err := service.taskGroup.Create(ctx, taskGroup)
	if err != nil {
		return nil, err
	}
	taskGroup.ID = id
	return taskGroup, nil
}

func (service *taskGroupService) Delete(ctx context.Context, userID uint64, groupID uint64) error {
	var err error = service.taskGroup.Delete(ctx, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskGroupService) Rename(ctx context.Context, userID uint64, groupID uint64, newTitle string) error {
	if newTitle == "" {
		return fmt.Errorf("title required")
	}
	var err error = service.taskGroup.UpdateTitle(ctx, groupID, userID, newTitle)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskGroupService) ListByUser(ctx context.Context, userID uint64) ([]models.TaskGroup, error) {
	taskGroups, err := service.taskGroup.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return taskGroups, nil
}

// Реализация бизнес-логики задачи

type taskService struct {
	task repository.TaskRepository
}

func (service *taskService) Create(ctx context.Context, userID uint64, groupId uint64, title string, body string) (*models.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	if body == "" {
		return nil, fmt.Errorf("body required")
	}
	var task *models.Task = &models.Task{
		UserID:  userID,
		GroupID: groupId,
		Title:   title,
		Body:    body,
	}
	id, err := service.task.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	task.ID = id
	return task, nil
}

func (service *taskService) Edit(ctx context.Context, userID uint64, taskID uint64, title string, body string) error {
	if title == "" {
		return fmt.Errorf("title required")
	}
	if body == "" {
		return fmt.Errorf("body required")
	}
	var err error = service.task.Update(ctx, taskID, userID, title, body)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskService) Delete(ctx context.Context, userID uint64, taskID uint64) error {
	var err error = service.task.Delete(ctx, taskID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskService) Move(ctx context.Context, userID uint64, taskID uint64, newGroupID uint64) error {
	var err error = service.task.MoveToGroup(ctx, taskID, userID, newGroupID)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskService) ToggleDone(ctx context.Context, userID uint64, taskID uint64, done bool) error {
	var err error = service.task.SetDone(ctx, taskID, userID, done)
	if err != nil {
		return err
	}
	return nil
}

func (service *taskService) ListByGroup(ctx context.Context, userId uint64, groupID uint64) ([]models.Task, error) {
	tasks, err := service.task.ListByGroup(ctx, userId, groupID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
