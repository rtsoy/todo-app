package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateTodoItemDTO struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

type TodoItem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed"`
}

type ListItem struct {
	ID     uuid.UUID
	ListID uuid.UUID
	ItemID uuid.UUID
}
