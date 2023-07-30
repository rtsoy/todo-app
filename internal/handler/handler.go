package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/service"
)

type Handler struct {
	*service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) InitRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := e.Group("/api", h.JWTAuthentication)
	{
		api.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello from protected route! :)")
		})
	}
}
