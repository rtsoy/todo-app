package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/service"
	mock_service "github.com/rtsoy/todo-app/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserServicer, input signInInput)

	tests := []struct {
		name                string
		inputBody           string
		inputData           signInInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@example.com", "password": "qwerty123"}`,
			inputData: signInInput{
				Email:    "test@example.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockUserServicer, input signInInput) {
				s.EXPECT().GenerateToken(input.Email, input.Password).Return("test-token", nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"token":"test-token"}`,
		},
		{
			name:                "Invalid JSON",
			inputBody:           `{`,
			mockBehavior:        func(s *mock_service.MockUserServicer, input signInInput) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"email": "test@example.com", "username": "test", "password": "qwerty123"}`,
			inputData: signInInput{
				Email:    "test@example.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockUserServicer, input signInInput) {
				s.EXPECT().GenerateToken(input.Email, input.Password).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUserServicer(c)
			test.mockBehavior(user, test.inputData)

			services := &service.Service{UserService: user}
			handler := NewHandler(services)

			e := echo.New()
			e.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody+"\n", w.Body.String())
		})
	}
}

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserServicer, user model.CreateUserDTO)

	tests := []struct {
		name                string
		inputBody           string
		inputUser           model.CreateUserDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@example.com", "username": "test", "password": "qwerty123"}`,
			inputUser: model.CreateUserDTO{
				Email:    "test@example.com",
				Username: "test",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockUserServicer, user model.CreateUserDTO) {
				s.EXPECT().CreateUser(user).Return(uuid.Nil, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"00000000-0000-0000-0000-000000000000"}`,
		},
		{
			name:                "Invalid JSON",
			inputBody:           `{`,
			mockBehavior:        func(s *mock_service.MockUserServicer, user model.CreateUserDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"Invalid JSON"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"email": "test@example.com", "username": "test", "password": "qwerty123"}`,
			inputUser: model.CreateUserDTO{
				Email:    "test@example.com",
				Username: "test",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockUserServicer, user model.CreateUserDTO) {
				s.EXPECT().CreateUser(user).Return(uuid.Nil, errors.New("service failure"))
			},
			expectedStatusCode:  http.StatusConflict,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUserServicer(c)
			test.mockBehavior(user, test.inputUser)

			services := &service.Service{UserService: user}
			handler := NewHandler(services)

			e := echo.New()
			e.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(test.inputBody))
			req.Header.Add("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody+"\n", w.Body.String())
		})
	}
}
