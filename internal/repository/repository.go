package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/internal/model"
)

type TodoItemRepository interface {
	Create(listID uuid.UUID, item model.CreateTodoItemDTO) (uuid.UUID, error)
	GetAll(userID, listID uuid.UUID) ([]model.TodoItem, error)
	GetByID(userID, itemID uuid.UUID) (model.TodoItem, error)
	Update(userID, itemID uuid.UUID, data model.UpdateTodoItemDTO) error
	Delete(userID, itemID uuid.UUID) error
}

type TodoListRepository interface {
	Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]model.TodoList, error)
	GetByID(userID, listID uuid.UUID) (model.TodoList, error)
	Update(userID, listID uuid.UUID, data model.UpdateTodoListDTO) error
	Delete(userID, listID uuid.UUID) error
}

type UserRepository interface {
	Create(user model.CreateUserDTO) (uuid.UUID, error)
	GetByEmail(email string) (*model.User, error)
}

type Repository struct {
	UserRepository
	TodoListRepository
	TodoItemRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     NewUserRepositoryPostgres(db),
		TodoListRepository: NewTodoListRepositoryPostgres(db),
		TodoItemRepository: NewTodoItemRepositoryPostgres(db),
	}
}
