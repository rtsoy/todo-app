package handler

import (
	"github.com/google/uuid"
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
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListByID)
			lists.PATCH("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)
		}
	}
}

func getIDFromParams(c echo.Context) (uuid.UUID, error) {
	paramsId := c.Param("id")

	return uuid.Parse(paramsId)
}

func getContextUserID(c echo.Context) uuid.UUID {
	ctxUserIDValue := c.Get(ctxUserID).(string)
	userID, _ := uuid.Parse(ctxUserIDValue)

	return userID
}
