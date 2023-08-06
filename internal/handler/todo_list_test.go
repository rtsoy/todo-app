package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/service"
	mock_service "github.com/rtsoy/todo-app/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_deleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID)

	tests := []struct {
		name                string
		listID              uuid.UUID
		listIDStr           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {
				s.EXPECT().Delete(userID, listID).Return(nil)
			},
			expectedStatusCode:  http.StatusNoContent,
			expectedRequestBody: "",
		},
		{
			name:                "Invalid ID",
			listID:              uuid.Nil,
			listIDStr:           "12312312",
			mockBehavior:        func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:      "Service Failure",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {
				s.EXPECT().Delete(userID, listID).Return(errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			userID := uuid.New()

			todoList := mock_service.NewMockTodoListServicer(c)
			test.mockBehavior(todoList, userID, test.listID)

			services := &service.Service{TodoListService: todoList}
			handler := NewHandler(services)

			e := echo.New()
			e.DELETE("/delete-list/:listID", handler.deleteList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/update-list/%s", test.listIDStr), nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("listID")
			ctx.SetParamValues(test.listIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.deleteList(ctx)
			if err != nil {
				httpErr := err.(*echo.HTTPError)

				errBytes, _ := json.Marshal(err)
				errJSON := string(errBytes)

				assert.Equal(t, test.expectedStatusCode, httpErr.Code)
				assert.Equal(t, test.expectedRequestBody, errJSON)
				return
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID, input model.UpdateTodoListDTO)

	tests := []struct {
		name                string
		listID              uuid.UUID
		listIDStr           string
		inputBody           string
		inputData           model.UpdateTodoListDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			inputBody: `{"title": "test", "description":"example"}`,
			inputData: model.UpdateTodoListDTO{},
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID, input model.UpdateTodoListDTO) {
				title := "test"
				description := "example"

				s.EXPECT().Update(userID, listID, model.UpdateTodoListDTO{
					Title:       &title,
					Description: &description,
				}).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ``,
		},
		{
			name:                "Invalid JSON",
			listID:              uuid.Nil,
			listIDStr:           uuid.Nil.String(),
			inputBody:           `{`,
			inputData:           model.UpdateTodoListDTO{},
			mockBehavior:        func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID, input model.UpdateTodoListDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:                "Invalid ID",
			listID:              uuid.Nil,
			listIDStr:           "12312312",
			mockBehavior:        func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID, input model.UpdateTodoListDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:      "Service Failure",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID, input model.UpdateTodoListDTO) {
				s.EXPECT().Update(userID, listID, input).Return(errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			userID := uuid.New()

			todoList := mock_service.NewMockTodoListServicer(c)
			test.mockBehavior(todoList, userID, test.listID, test.inputData)

			services := &service.Service{TodoListService: todoList}
			handler := NewHandler(services)

			e := echo.New()
			e.PATCH("/update-list/:listID", handler.updateList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/update-list/%s", test.listIDStr), bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("listID")
			ctx.SetParamValues(test.listIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.updateList(ctx)
			if err != nil {
				httpErr := err.(*echo.HTTPError)

				errBytes, _ := json.Marshal(err)
				errJSON := string(errBytes)

				assert.Equal(t, test.expectedStatusCode, httpErr.Code)
				assert.Equal(t, test.expectedRequestBody, errJSON)
				return
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getListByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID)

	tests := []struct {
		name                string
		listID              uuid.UUID
		listIDStr           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {
				s.EXPECT().GetByID(userID, listID).Return(model.TodoList{
					ID:          listID,
					Title:       "test",
					Description: "example",
					CreatedAt:   time.Unix(0, 0),
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"00000000-0000-0000-0000-000000000000","title":"test","description":"example","createdAt":"1970-01-01T06:00:00+06:00"}`,
		},
		{
			name:                "Invalid ID",
			listID:              uuid.Nil,
			listIDStr:           "12312312",
			mockBehavior:        func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:      "Service Failure",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID, listID uuid.UUID) {
				s.EXPECT().GetByID(userID, listID).Return(model.TodoList{}, errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			userID := uuid.New()

			todoList := mock_service.NewMockTodoListServicer(c)
			test.mockBehavior(todoList, userID, test.listID)

			services := &service.Service{TodoListService: todoList}
			handler := NewHandler(services)

			e := echo.New()
			e.GET("/get-list-by-id/:listID", handler.getListByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(
				"/get-list-by-id/%s", test.listIDStr), nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("listID")
			ctx.SetParamValues(test.listIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.getListByID(ctx)
			if err != nil {
				httpErr := err.(*echo.HTTPError)

				errBytes, _ := json.Marshal(err)
				errJSON := string(errBytes)

				assert.Equal(t, test.expectedStatusCode, httpErr.Code)
				assert.Equal(t, test.expectedRequestBody, errJSON)
				return
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody+"\n", w.Body.String())
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoListServicer, userID uuid.UUID)

	tests := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID uuid.UUID) {
				s.EXPECT().GetAll(userID, nil).Return([]model.TodoList{
					{
						ID:          uuid.Nil,
						Title:       "test1",
						Description: "example",
						CreatedAt:   time.Unix(0, 0),
					},
					{
						ID:          uuid.Nil,
						Title:       "test2",
						Description: "example",
						CreatedAt:   time.Unix(0, 0),
					},
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"count":2,"results":[{"id":"00000000-0000-0000-0000-000000000000","title":"test1","description":"example","createdAt":"1970-01-01T06:00:00+06:00"},{"id":"00000000-0000-0000-0000-000000000000","title":"test2","description":"example","createdAt":"1970-01-01T06:00:00+06:00"}],"pagination":null}`,
		},
		{
			name: "Service Failure",
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID uuid.UUID) {
				s.EXPECT().GetAll(userID, nil).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uuid.New()

			todoList := mock_service.NewMockTodoListServicer(c)
			test.mockBehavior(todoList, userID)

			services := &service.Service{TodoListService: todoList}
			handler := NewHandler(services)

			e := echo.New()
			e.GET("/get-all-lists", handler.getAllLists)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/get-all-lists", nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)
			ctx.Set(ctxUserID, userID.String())
			err := handler.getAllLists(ctx)
			if err != nil {
				httpErr := err.(*echo.HTTPError)

				errBytes, _ := json.Marshal(err)
				errJSON := string(errBytes)

				assert.Equal(t, test.expectedStatusCode, httpErr.Code)
				assert.Equal(t, test.expectedRequestBody, errJSON)
				return
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody+"\n", w.Body.String())
		})
	}
}

func TestHandler_createList(t *testing.T) {

	type mockBehavior func(s *mock_service.MockTodoListServicer, userID uuid.UUID, input model.CreateTodoListDTO)

	tests := []struct {
		name                string
		inputBody           string
		inputData           model.CreateTodoListDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"test", "description":"example"}`,
			inputData: model.CreateTodoListDTO{
				Title:       "test",
				Description: "example",
			},
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID uuid.UUID, input model.CreateTodoListDTO) {
				s.EXPECT().Create(userID, input).Return(uuid.Nil, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: `{"id":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			name:                "Invalid JSON",
			inputBody:           `{`,
			mockBehavior:        func(s *mock_service.MockTodoListServicer, userID uuid.UUID, input model.CreateTodoListDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"title":"test", "description":"example"}`,
			inputData: model.CreateTodoListDTO{
				Title:       "test",
				Description: "example",
			},
			mockBehavior: func(s *mock_service.MockTodoListServicer, userID uuid.UUID, input model.CreateTodoListDTO) {
				s.EXPECT().Create(userID, input).Return(uuid.Nil, errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uuid.New()

			todoList := mock_service.NewMockTodoListServicer(c)
			test.mockBehavior(todoList, userID, test.inputData)

			services := &service.Service{TodoListService: todoList}
			handler := NewHandler(services)

			e := echo.New()
			e.POST("/create-list", handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/create-list", bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)
			ctx.Set(ctxUserID, userID.String())
			err := handler.createList(ctx)
			if err != nil {
				httpErr := err.(*echo.HTTPError)

				errBytes, _ := json.Marshal(err)
				errJSON := string(errBytes)

				assert.Equal(t, test.expectedStatusCode, httpErr.Code)
				assert.Equal(t, test.expectedRequestBody, errJSON)
				return
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody+"\n", w.Body.String())
		})
	}
}
