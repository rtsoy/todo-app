package service

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
)

const (
	minListTitleLength       = 3
	minListDescriptionLength = 3
	maxListsPerUser          = 5
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
	totalLists, err := s.repository.GetAll(userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, err
	}

	if len(totalLists) >= maxListsPerUser {
		return uuid.Nil, errors.New("exceeded the maximum allowed limit of existing lists")
	}

	if len(list.Title) < minListTitleLength {
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

	if lists == nil {
		return lists, errors.New("no todo lists found")
	}

	return lists, nil
}

func (s *TodoListService) GetByID(userID, listID uuid.UUID) (model.TodoList, error) {
	list, err := s.repository.GetByID(userID, listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return list, errors.New("todo list not found")
		}

		return list, err
	}

	return list, nil
}

func (s *TodoListService) Update(userID, listID uuid.UUID, data model.UpdateTodoListDTO) error {
	if reflect.DeepEqual(data, model.UpdateTodoListDTO{}) {
		return errors.New("there is no values to update")
	}

	if data.Title != nil && len(*data.Title) < minListTitleLength && len(*data.Title) > 0 {
		return errors.New("title length is too short")
	}

	if data.Description != nil && len(*data.Description) < minListDescriptionLength && len(*data.Description) > 0 {
		return errors.New("description length is too short")
	}

	return s.repository.Update(userID, listID, data)
}

func (s *TodoListService) Delete(userID, listID uuid.UUID) error {
	return s.repository.Delete(userID, listID)
}
