package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/internal/model"
)

type TodoListRepository interface {
	Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]model.TodoList, error)
	GetByID(userID, listID uuid.UUID) (model.TodoList, error)
	Update(userID, listID uuid.UUID, data model.CreateTodoListDTO) error
	Delete(userID, listID uuid.UUID) error
}

type UserRepository interface {
	Create(user model.CreateUserDTO) (uuid.UUID, error)
	GetByEmail(email string) (*model.User, error)
}

type Repository struct {
	UserRepository
	TodoListRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     NewUserRepositoryPostgres(db),
		TodoListRepository: NewTodoListRepositoryPostgres(db),
	}
}
