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
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type TodoItemRepositoryPostgres struct {
	db *sqlx.DB
}

func NewTodoItemRepositoryPostgres(db *sqlx.DB) TodoItemRepository {
	return &TodoItemRepositoryPostgres{
		db: db,
	}
}

func (r *TodoItemRepositoryPostgres) Create(listID uuid.UUID, item model.CreateTodoItemDTO) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	createItemQuery := fmt.Sprintf(`
		INSERT INTO %s (id, title, description, created_at, deadline, completed)
		VALUES ($1, $2, $3, $4, $5, $6)
    `, todoItemsTable)

	itemID := uuid.New()
	if _, err := tx.Exec(createItemQuery, itemID, item.Title, item.Description, time.Now().UTC(), item.Deadline, false); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	createListItemQuery := fmt.Sprintf(`
		INSERT INTO %s (id, list_id, item_id)
		VALUES ($1, $2, $3)
    `, listsItemsTable)

	if _, err := tx.Exec(createListItemQuery, uuid.New(), listID, itemID); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return itemID, tx.Commit()
}

func (r *TodoItemRepositoryPostgres) GetAll(userID, listID uuid.UUID) ([]model.TodoItem, error) {
	query := fmt.Sprintf(`
		SELECT ti.id, ti.title, ti.description, ti.created_at, ti.deadline, ti.completed
		FROM %s ti
		INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id
		WHERE ul.user_id = $1 AND li.list_id = $2
    `, todoItemsTable, listsItemsTable, usersListsTable)

	var items []model.TodoItem

	return items, r.db.Select(&items, query, userID, listID)
}

func (r *TodoItemRepositoryPostgres) GetByID(userID, itemID uuid.UUID) (model.TodoItem, error) {
	query := fmt.Sprintf(`
		SELECT ti.id, ti.title, ti.description, ti.created_at, ti.deadline, ti.completed
		FROM %s ti
		INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id
		WHERE ul.user_id = $1 AND ti.id = $2
    `, todoItemsTable, listsItemsTable, usersListsTable)
	var item model.TodoItem

	return item, r.db.Get(&item, query, userID, itemID)
}

func (r *TodoItemRepositoryPostgres) Update(userID, itemID uuid.UUID, data model.CreateTodoItemDTO) error {
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

	if !data.Deadline.IsZero() {
		toUpdate = append(toUpdate, fmt.Sprintf("deadline=$%d", argsID))
		args = append(args, data.Deadline)
		argsID++
	}

	updateQuery := strings.Join(toUpdate, ", ")
	args = append(args, userID, itemID)

	query := fmt.Sprintf(`
		UPDATE %s ti
		SET %s
		FROM %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d
    `, todoItemsTable, updateQuery, listsItemsTable, usersListsTable, argsID, argsID+1)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoItemRepositoryPostgres) Delete(userID, itemID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s ti
		USING %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2
    `, todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userID, itemID)

	return err
}
