package service

import (
	"database/sql"
	"errors"
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

	if time.Now().UTC().After(item.Deadline.UTC()) {
		return uuid.Nil, errors.New("deadline cannot be in the past")
	}

	return s.repository.Create(listID, item)
}

func (s *TodoItemService) GetAll(userID, listID uuid.UUID) ([]model.TodoItem, error) {
	items, err := s.repository.GetAll(userID, listID)
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

func (s *TodoItemService) Update(userID, itemID uuid.UUID, data model.CreateTodoItemDTO) error {
	if data.Title == "" && data.Description == "" && data.Deadline.IsZero() {
		return errors.New("title, description and deadline cannot be empty")
	}

	if len(data.Title) < minItemTitleLength && len(data.Title) > 0 {
		return errors.New("title length is too short")
	}

	if len(data.Description) < minItemDescriptionLength && len(data.Description) > 0 {
		return errors.New("description length is too short")
	}

	item, _ := s.repository.GetByID(userID, itemID)

	if item.CreatedAt.After(data.Deadline) {
		return errors.New("deadline cannot be in the past")
	}

	return s.repository.Update(userID, itemID, data)
}

func (s *TodoItemService) Delete(userID, itemID uuid.UUID) error {
	return s.repository.Delete(userID, itemID)
}
