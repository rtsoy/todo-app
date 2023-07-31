package service

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
)

const (
	minListTileLength        = 3
	minListDescriptionLength = 3
)

type TodoListService struct {
	repository repository.TodoListRepository
}

func NewTodoListService(repository repository.TodoListRepository) TodoListServicer {
	return &TodoListService{
		repository: repository,
	}
}

func (s *TodoListService) Create(userID uuid.UUID, list model.CreateTodoListDTO) (uuid.UUID, error) {
	if len(list.Title) < minListTileLength {
		return uuid.Nil, errors.New("title length is too short")
	}

	if len(list.Description) < minListDescriptionLength {
		return uuid.Nil, errors.New("description length is too short")
	}

	return s.repository.Create(userID, list)
}

func (s *TodoListService) GetAll(userID uuid.UUID) ([]model.TodoList, error) {
	lists, err := s.repository.GetAll(userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lists, errors.New("no todo lists found")
		}

		return lists, err
	}

	return lists, err
}

func (s *TodoListService) GetByID(userID, listID uuid.UUID) (model.TodoList, error) {
	list, err := s.repository.GetByID(userID, listID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return list, errors.New("todo list not found")
		}

		return list, err
	}

	return list, err
}

func (s *TodoListService) Update(userID, listID uuid.UUID, data model.CreateTodoListDTO) error {
	if data.Title == "" && data.Description == "" {
		return errors.New("title and description cannot be empty")
	}

	return s.repository.Update(userID, listID, data)
}

func (s *TodoListService) Delete(userID, listID uuid.UUID) error {
	return s.repository.Delete(userID, listID)
}
