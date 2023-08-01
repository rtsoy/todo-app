package model

import (
	"time"

	"github.com/google/uuid"
)

type UpdateTodoListDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CreateTodoListDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoList struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
