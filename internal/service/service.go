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

type Service struct {
	UserService UserServicer
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repository.UserRepository),
	}
}
