package service

import (
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
)

const (
	minItemTitleLength       = 3
	minItemDescriptionLength = 3
)

type TodoItemService struct {
	repository     repository.TodoItemRepository
	listRepository repository.TodoListRepository
}

func NewTodoItemService(repository repository.TodoItemRepository, listRepository repository.TodoListRepository) TodoItemServicer {
	return &TodoItemService{
		repository:     repository,
		listRepository: listRepository,
	}
}

func (s *TodoItemService) Create(userID, listID uuid.UUID, item model.CreateTodoItemDTO) (uuid.UUID, error) {
	if _, err := s.listRepository.GetByID(userID, listID); err != nil {
		return uuid.Nil, errors.New("forbidden")
	}

	if len(item.Title) < minItemTitleLength {
		return uuid.Nil, errors.New("title length is too short")
	}

	if len(item.Description) < minItemDescriptionLength {
		return uuid.Nil, errors.New("description length is too short")
	}

	// Default deadline is 7 days
	if item.Deadline.IsZero() {
		item.Deadline = time.Now().UTC().AddDate(0, 0, 7)
	}

	if time.Now().UTC().After(item.Deadline.UTC()) {
		return uuid.Nil, errors.New("deadline cannot be in the past")
	}

	return s.repository.Create(listID, item)
}

func (s *TodoItemService) GetAll(userID, listID uuid.UUID, pagination *model.Pagination) ([]model.TodoItem, error) {
	if pagination.Limit == 0 {
		pagination.Limit = 5
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}

	items, err := s.repository.GetAll(userID, listID, pagination)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return items, errors.New("no todo items found")
		}

		return items, err
	}

	if items == nil {
		return items, errors.New("no todo items found")
	}

	return items, nil
}

func (s *TodoItemService) GetByID(userID, itemID uuid.UUID) (model.TodoItem, error) {
	list, err := s.repository.GetByID(userID, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return list, errors.New("todo list not found")
		}

		return list, err
	}

	return list, nil
}

func (s *TodoItemService) Update(userID, itemID uuid.UUID, data model.UpdateTodoItemDTO) error {
	if reflect.DeepEqual(data, model.UpdateTodoItemDTO{}) {
		return errors.New("there is no values to update")
	}

	if data.Title != nil && len(*data.Title) < minItemTitleLength && len(*data.Title) > 0 {
		return errors.New("title length is too short")
	}

	if data.Description != nil && len(*data.Description) < minItemDescriptionLength && len(*data.Description) > 0 {
		return errors.New("description length is too short")
	}

	item, _ := s.repository.GetByID(userID, itemID)
	if data.Deadline != nil && item.CreatedAt.After(*data.Deadline) {
		return errors.New("deadline cannot be in the past")
	}

	return s.repository.Update(userID, itemID, data)
}

func (s *TodoItemService) Delete(userID, itemID uuid.UUID) error {
	return s.repository.Delete(userID, itemID)
}
