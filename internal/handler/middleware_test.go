package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/service"
	mock_service "github.com/rtsoy/todo-app/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_JWTAuthentication(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserServicer, token string)

	tests := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUserServicer, token string) {
				s.EXPECT().ParseToken(token).Return(jwt.MapClaims{
					"userID":  1,
					"expires": time.Now().Add(time.Minute).Format(time.RFC3339),
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"userID":1}`,
		},
		{
			name:                "Empty Auth Header",
			headerName:          "Authorization",
			headerValue:         "",
			token:               "",
			mockBehavior:        func(s *mock_service.MockUserServicer, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"empty auth header"}`,
		},
		{
			name:                "Invalid Auth Header",
			headerName:          "Authorization",
			headerValue:         "Baerer",
			token:               "",
			mockBehavior:        func(s *mock_service.MockUserServicer, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"invalid auth token"}`,
		},
		{
			name:                "Invalid Auth Header",
			headerName:          "Authorization",
			headerValue:         "Bearer qwe 1",
			token:               "",
			mockBehavior:        func(s *mock_service.MockUserServicer, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"invalid auth token"}`,
		},
		{
			name:                "Invalid Auth Header",
			headerName:          "Authorization",
			headerValue:         "Bearer ",
			token:               "",
			mockBehavior:        func(s *mock_service.MockUserServicer, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"no token provided"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUserServicer, token string) {
				s.EXPECT().ParseToken(token).Return(nil, errors.New("parse error"))
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"parse error"}`,
		},
		{
			name:        "Expired Token",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockUserServicer, token string) {
				s.EXPECT().ParseToken(token).Return(jwt.MapClaims{
					"userID":  1,
					"expires": time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
				}, nil)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"token is expired"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUserServicer(c)
			test.mockBehavior(user, test.token)

			services := &service.Service{UserService: user}
			handler := NewHandler(services)

			e := echo.New()
			e.GET("/jwt", handler.JWTAuthentication(func(c echo.Context) error {
				userID := c.Get(ctxUserID)

				return c.JSON(http.StatusOK, echo.Map{
					"userID": userID,
				})
			}))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/jwt", nil)
			req.Header.Add("Authorization", test.headerValue)

			e.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedRequestBody+"\n")
		})
	}
}
