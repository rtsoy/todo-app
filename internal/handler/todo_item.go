package handler

import (
	"fmt"
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

	var input model.UpdateTodoItemDTO
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

	var pagination model.Pagination
	if err := c.Bind(&pagination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid url query")
	}

	fmt.Println(c.QueryParams())

	items, err := h.TodoItemService.GetAll(userID, listID, &pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, resourceResponse{
		Count:      len(items),
		Results:    items,
		Pagination: &pagination,
	})
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
