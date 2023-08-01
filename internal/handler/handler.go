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
			lists.GET("/:listID", h.getListByID)
			lists.PATCH("/:listID", h.updateList)
			lists.DELETE("/:listID", h.deleteList)

			items := lists.Group("/:listID/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:itemID", h.getItemByID)
				items.PATCH("/:itemID", h.updateItem)
				items.DELETE("/:itemID", h.deleteItem)
			}
		}
	}
}

func getValueFromParams(c echo.Context, v string) (uuid.UUID, error) {
	paramsId := c.Param(v)

	return uuid.Parse(paramsId)
}

func getContextUserID(c echo.Context) uuid.UUID {
	ctxUserIDValue := c.Get(ctxUserID).(string)
	userID, _ := uuid.Parse(ctxUserIDValue)

	return userID
}
