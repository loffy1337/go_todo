package models

import "time"

// Типы для группы задач
const (
	GroupStatusUrgent   = "urgent"
	GroupStatusDaily    = "daily"
	GroupStatusLongTerm = "longterm"
)

// Модель пользователя
type User struct {
	ID           uint64     `json:"id"`
	Username     string     `json:"username"`
	Password     string     `json:"-"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	StatusText   string     `json:"status_text"`
	RegisteredAt time.Time  `json:"registered_at"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	IsOnline     bool       `json:"is_online"`
}

// Модель группы задач
type TaskGroup struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Модель задачи
type Task struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	GroupID   uint64    `json:"group_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
