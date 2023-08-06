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

func TestHandler_deleteItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID)

	tests := []struct {
		name                string
		itemID              uuid.UUID
		itemIDStr           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {
				s.EXPECT().Delete(userID, itemID).Return(nil)
			},
			expectedStatusCode:  http.StatusNoContent,
			expectedRequestBody: "",
		},
		{
			name:                "Invalid ID",
			itemID:              uuid.Nil,
			itemIDStr:           "123123123",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:      "Service Failure",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {
				s.EXPECT().Delete(userID, itemID).Return(errors.New("service failure"))
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

			todoItem := mock_service.NewMockTodoItemServicer(c)
			test.mockBehavior(todoItem, userID, test.itemID)

			services := &service.Service{TodoItemService: todoItem}
			handler := NewHandler(services)

			e := echo.New()
			e.DELETE("/delete-item/:id", handler.deleteItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/delete-item/%s", test.itemIDStr), nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("itemID")
			ctx.SetParamValues(test.itemIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.deleteItem(ctx)
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

func TestHandler_updateItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID, input model.UpdateTodoItemDTO)

	tests := []struct {
		name                string
		itemID              uuid.UUID
		itemIDStr           string
		inputBody           string
		inputData           model.UpdateTodoItemDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			inputBody: `{"title":"test", "description":"example", "deadline":"1970-01-01T00:00:00Z", "completed":true}`,
			inputData: model.UpdateTodoItemDTO{},
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID, input model.UpdateTodoItemDTO) {
				title := "test"
				description := "example"
				deadline := time.Unix(0, 0).UTC()
				completed := true

				s.EXPECT().Update(userID, itemID, model.UpdateTodoItemDTO{
					Title:       &title,
					Description: &description,
					Deadline:    &deadline,
					Completed:   &completed,
				}).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ``,
		},
		{
			name:                "Invalid ID",
			itemID:              uuid.Nil,
			itemIDStr:           "123123123",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID, input model.UpdateTodoItemDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:                "Invalid JSON",
			itemID:              uuid.Nil,
			itemIDStr:           uuid.Nil.String(),
			inputBody:           `{`,
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID, input model.UpdateTodoItemDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:      "Service Failure",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID, input model.UpdateTodoItemDTO) {
				s.EXPECT().Update(userID, itemID, input).Return(errors.New("service failure"))
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

			todoItem := mock_service.NewMockTodoItemServicer(c)
			test.mockBehavior(todoItem, userID, test.itemID, test.inputData)

			services := &service.Service{TodoItemService: todoItem}
			handler := NewHandler(services)

			e := echo.New()
			e.PATCH("/update-item/:id", handler.updateItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/update-item/%s", test.itemIDStr), bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("itemID")
			ctx.SetParamValues(test.itemIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.updateItem(ctx)
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

func TestHandler_getItemByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID)

	tests := []struct {
		name                string
		itemID              uuid.UUID
		itemIDStr           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {
				s.EXPECT().GetByID(userID, itemID).Return(model.TodoItem{
					ID:          uuid.Nil,
					Title:       "test",
					Description: "example",
					CreatedAt:   time.Unix(0, 0),
					Deadline:    time.Unix(0, 1),
					Completed:   true,
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"00000000-0000-0000-0000-000000000000","title":"test","description":"example","createdAt":"1970-01-01T06:00:00+06:00","deadline":"1970-01-01T06:00:00.000000001+06:00","completed":true}`,
		},
		{
			name:                "Invalid ID",
			itemID:              uuid.Nil,
			itemIDStr:           "12312321",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid id"}`,
		},
		{
			name:      "Service Failure",
			itemID:    uuid.Nil,
			itemIDStr: uuid.Nil.String(),
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, itemID uuid.UUID) {
				s.EXPECT().GetByID(userID, itemID).Return(model.TodoItem{}, errors.New("service failure"))
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

			todoItem := mock_service.NewMockTodoItemServicer(c)
			test.mockBehavior(todoItem, userID, test.itemID)

			services := &service.Service{TodoItemService: todoItem}
			handler := NewHandler(services)

			e := echo.New()
			e.GET("/get-item-by-id/:id", handler.getItemByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get-item-by-id/%s", test.itemIDStr), nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)

			ctx.SetParamNames("itemID")
			ctx.SetParamValues(test.itemIDStr)

			ctx.Set(ctxUserID, userID.String())
			err := handler.getItemByID(ctx)
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

func TestHandler_getAllItems(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID)

	tests := []struct {
		name                string
		listID              uuid.UUID
		listIDStr           string
		paginationPage      string
		paginationLimit     string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:            "OK",
			listID:          uuid.Nil,
			listIDStr:       uuid.Nil.String(),
			paginationPage:  "1",
			paginationLimit: "5",
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID) {
				s.EXPECT().GetAll(userID, listID, &model.Pagination{Page: 1, Limit: 5}, nil).Return([]model.TodoItem{
					{
						ID:          uuid.Nil,
						Title:       "test1",
						Description: "example",
						CreatedAt:   time.Unix(0, 0),
						Deadline:    time.Unix(0, 1),
						Completed:   false,
					},
					{
						ID:          uuid.Nil,
						Title:       "test2",
						Description: "example",
						CreatedAt:   time.Unix(0, 0),
						Deadline:    time.Unix(0, 1),
						Completed:   true,
					},
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"count":2,"results":[{"id":"00000000-0000-0000-0000-000000000000","title":"test1","description":"example","createdAt":"1970-01-01T06:00:00+06:00","deadline":"1970-01-01T06:00:00.000000001+06:00","completed":false},{"id":"00000000-0000-0000-0000-000000000000","title":"test2","description":"example","createdAt":"1970-01-01T06:00:00+06:00","deadline":"1970-01-01T06:00:00.000000001+06:00","completed":true}],"pagination":{"page":1,"limit":5}}`,
		},
		{
			name:                "Invalid ListID",
			paginationPage:      "1",
			paginationLimit:     "5",
			listID:              uuid.Nil,
			listIDStr:           "12312321",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid list id"}`,
		},
		{
			name:                "Invalid URL Query",
			listID:              uuid.Nil,
			listIDStr:           uuid.Nil.String(),
			paginationPage:      "qweqwe",
			paginationLimit:     "###",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid url query"}`,
		},
		{
			name:            "Service Failure",
			listID:          uuid.Nil,
			listIDStr:       uuid.Nil.String(),
			paginationPage:  "1",
			paginationLimit: "5",
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID) {
				s.EXPECT().GetAll(userID, listID, &model.Pagination{Page: 1, Limit: 5}, nil).Return(nil, errors.New("service failure"))
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

			todoItem := mock_service.NewMockTodoItemServicer(c)
			test.mockBehavior(todoItem, userID, test.listID)

			services := &service.Service{TodoItemService: todoItem}
			handler := NewHandler(services)

			e := echo.New()
			e.GET("/get-all-items", handler.getAllLists)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/get-all-items", nil)
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)
			ctx.Set(ctxUserID, userID.String())

			ctx.QueryParams().Add("page", test.paginationPage)
			ctx.QueryParams().Add("limit", test.paginationLimit)

			ctx.SetParamNames("listID")
			ctx.SetParamValues(test.listIDStr)

			err := handler.getAllItems(ctx)
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

func TestHandler_createItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID, input model.CreateTodoItemDTO)

	tests := []struct {
		name                string
		listID              uuid.UUID
		listIDStr           string
		inputBody           string
		inputData           model.CreateTodoItemDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			inputBody: `{"title":"test", "description":"example", "deadline":"1970-01-01T00:00:00Z"}`,
			inputData: model.CreateTodoItemDTO{
				Title:       "test",
				Description: "example",
				Deadline:    time.Unix(0, 0).UTC(),
			},
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID, input model.CreateTodoItemDTO) {
				s.EXPECT().Create(userID, listID, input).Return(uuid.Nil, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: `{"id":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			name:                "Invalid ListID",
			listID:              uuid.Nil,
			listIDStr:           "12312312",
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID, input model.CreateTodoItemDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid list id"}`,
		},
		{
			name:                "Invalid JSON",
			listID:              uuid.Nil,
			listIDStr:           uuid.Nil.String(),
			inputBody:           `{`,
			mockBehavior:        func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID, input model.CreateTodoItemDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:      "Service Failure",
			listID:    uuid.Nil,
			listIDStr: uuid.Nil.String(),
			inputBody: `{"title":"test", "description":"example", "deadline":"1970-01-01T00:00:00Z"}`,
			inputData: model.CreateTodoItemDTO{
				Title:       "test",
				Description: "example",
				Deadline:    time.Unix(0, 0).UTC(),
			},
			mockBehavior: func(s *mock_service.MockTodoItemServicer, userID, listID uuid.UUID, input model.CreateTodoItemDTO) {
				s.EXPECT().Create(userID, listID, input).Return(uuid.Nil, errors.New("service failure"))
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

			todoItem := mock_service.NewMockTodoItemServicer(c)
			test.mockBehavior(todoItem, userID, test.listID, test.inputData)

			services := &service.Service{TodoItemService: todoItem}
			handler := NewHandler(services)

			e := echo.New()
			e.POST("/create-item", handler.createItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/create-item", bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			ctx := e.NewContext(req, w)
			ctx.Set(ctxUserID, userID.String())

			ctx.SetParamNames("listID")
			ctx.SetParamValues(test.listIDStr)

			err := handler.createItem(ctx)
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
