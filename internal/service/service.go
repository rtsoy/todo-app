package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
)

type TodoItemServicer interface {
	Create(userID, listID uuid.UUID, item model.CreateTodoItemDTO) (uuid.UUID, error)
	GetAll(userID, listID uuid.UUID) ([]model.TodoItem, error)
	GetByID(userID, itemID uuid.UUID) (model.TodoItem, error)
	Update(userID, itemID uuid.UUID, data model.UpdateTodoItemDTO) error
	Delete(userID, itemID uuid.UUID) error
}

type TodoListServicer interface {
	Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]model.TodoList, error)
	GetByID(userID, listID uuid.UUID) (model.TodoList, error)
	Update(userID, listID uuid.UUID, data model.UpdateTodoListDTO) error
	Delete(userID, listID uuid.UUID) error
}

type UserServicer interface {
	CreateUser(user model.CreateUserDTO) (uuid.UUID, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (jwt.MapClaims, error)
}

type Service struct {
	UserService     UserServicer
	TodoListService TodoListServicer
	TodoItemService TodoItemServicer
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoItemService: NewTodoItemService(repository.TodoItemRepository, repository.TodoListRepository),
		TodoListService: NewTodoListService(repository.TodoListRepository),
		UserService:     NewUserService(repository.UserRepository),
	}
}
