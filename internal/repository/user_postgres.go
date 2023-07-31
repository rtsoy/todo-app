package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/internal/model"
)

const usersTable = "users"

type UserRepositoryPostgres struct {
	db *sqlx.DB
}

func NewUserRepositoryPostgres(db *sqlx.DB) UserRepository {
	return &UserRepositoryPostgres{
		db: db,
	}
}

func (r *UserRepositoryPostgres) GetByEmail(email string) (*model.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, username, password_hash
		FROM %s
		WHERE email = $1
	`, usersTable)

	var user model.User

	return &user, r.db.Get(&user, query, email)
}

func (r *UserRepositoryPostgres) Create(user model.CreateUserDTO) (uuid.UUID, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, email, username, password_hash)
		VALUES ($1, $2, $3, $4)
	`, usersTable)

	id := uuid.New()

	return id, r.db.QueryRow(query, id, user.Email, user.Username, user.Password).Err()
}
