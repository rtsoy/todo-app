package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/internal/model"
)

const (
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
)

type TodoListRepositoryPostgres struct {
	db *sqlx.DB
}

func NewTodoListRepositoryPostgres(db *sqlx.DB) TodoListRepository {
	return &TodoListRepositoryPostgres{
		db: db,
	}
}

func (r *TodoListRepositoryPostgres) Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	createListQuery := fmt.Sprintf(`
		INSERT INTO %s (id, title, description, created_at)
		VALUES ($1, $2, $3, $4)
    `, todoListsTable)

	listID := uuid.New()
	if _, err := tx.Exec(createListQuery, listID, list.Title, list.Description, time.Now().UTC()); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	createUserListQuery := fmt.Sprintf(`
		INSERT INTO %s (id, user_id, list_id)
		VALUES ($1, $2, $3)
    `, usersListsTable)

	if _, err := tx.Exec(createUserListQuery, uuid.New(), userID, listID); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return listID, tx.Commit()
}

func (r *TodoListRepositoryPostgres) GetAll(userID uuid.UUID) ([]model.TodoList, error) {
	query := fmt.Sprintf(`
		SELECT tl.id, tl.title, tl.description, tl.created_at
		FROM %s tl
		INNER JOIN %s ul ON tl.id = ul.list_id
		WHERE ul.user_id = $1
    `, todoListsTable, usersListsTable)

	var lists []model.TodoList

	return lists, r.db.Select(&lists, query, userID)
}

func (r *TodoListRepositoryPostgres) GetByID(userID, listID uuid.UUID) (model.TodoList, error) {
	query := fmt.Sprintf(`
		SELECT tl.id, tl.title, tl.description, tl.created_at
		FROM %s tl
		INNER JOIN %s ul ON tl.id = ul.list_id
		WHERE ul.user_id = $1 AND ul.list_id = $2
    `, todoListsTable, usersListsTable)

	var list model.TodoList

	return list, r.db.Get(&list, query, userID, listID)
}

func (r *TodoListRepositoryPostgres) Update(userID, listID uuid.UUID, data model.CreateTodoListDTO) error {
	toUpdate := make([]string, 0)

	args := make([]interface{}, 0)
	argsID := 1

	if data.Title != "" {
		toUpdate = append(toUpdate, fmt.Sprintf("title=$%d", argsID))
		args = append(args, data.Title)
		argsID++
	}

	if data.Description != "" {
		toUpdate = append(toUpdate, fmt.Sprintf("description=$%d", argsID))
		args = append(args, data.Description)
		argsID++
	}

	updateQuery := strings.Join(toUpdate, ", ")
	args = append(args, listID, userID)

	query := fmt.Sprintf(`
		UPDATE %s tl
		SET %s
		FROM %s ul
		WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d
    `, todoListsTable, updateQuery, usersListsTable, argsID, argsID+1)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoListRepositoryPostgres) Delete(userID, listID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s tl
		USING %s ul
		WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2
    `, todoListsTable, usersListsTable)

	_, err := r.db.Exec(query, userID, listID)

	return err
}
