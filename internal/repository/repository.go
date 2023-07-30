package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/internal/model"
)

type UserRepository interface {
	Create(user model.CreateUserDTO) (uuid.UUID, error)
	GetByEmail(email string) (*model.User, error)
}

type Repository struct {
	UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepositoryPostgres(db),
	}
}
