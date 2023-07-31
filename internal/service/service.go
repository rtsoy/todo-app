package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
)

type UserServicer interface {
	CreateUser(user model.CreateUserDTO) (uuid.UUID, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (jwt.MapClaims, error)
}

type TodoListServicer interface {
	Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]model.TodoList, error)
	GetByID(userID, listID uuid.UUID) (model.TodoList, error)
	Update(userID, listID uuid.UUID, data model.CreateTodoListDTO) error
	Delete(userID, listID uuid.UUID) error
}

type Service struct {
	UserService     UserServicer
	TodoListService TodoListServicer
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoListService: NewTodoListService(repository.TodoListRepository),
		UserService:     NewUserService(repository.UserRepository),
	}
}
