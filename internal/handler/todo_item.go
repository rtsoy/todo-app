package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
)

func (h *Handler) deleteItem(c echo.Context) error {
	userID := getContextUserID(c)

	itemID, err := getValueFromParams(c, "itemID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.TodoItemService.Delete(userID, itemID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) updateItem(c echo.Context) error {
	userID := getContextUserID(c)

	itemID, err := getValueFromParams(c, "itemID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var input model.CreateTodoItemDTO
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := h.TodoItemService.Update(userID, itemID, input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) getItemByID(c echo.Context) error {
	userID := getContextUserID(c)

	itemID, err := getValueFromParams(c, "itemID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	item, err := h.TodoItemService.GetByID(userID, itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}

func (h *Handler) getAllItems(c echo.Context) error {
	userID := getContextUserID(c)

	listID, err := getValueFromParams(c, "listID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid list id")
	}

	items, err := h.TodoItemService.GetAll(userID, listID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

func (h *Handler) createItem(c echo.Context) error {
	userID := getContextUserID(c)

	listID, err := getValueFromParams(c, "listID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid list id")
	}

	var input model.CreateTodoItemDTO
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	id, err := h.TodoItemService.Create(userID, listID, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id,
	})
}
