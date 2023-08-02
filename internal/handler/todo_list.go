package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
)

func (h *Handler) deleteList(c echo.Context) error {
	userID := getContextUserID(c)

	listID, err := getValueFromParams(c, "listID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.TodoListService.Delete(userID, listID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) updateList(c echo.Context) error {
	userID := getContextUserID(c)

	listID, err := getValueFromParams(c, "listID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var input model.UpdateTodoListDTO
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := h.TodoListService.Update(userID, listID, input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) getListByID(c echo.Context) error {
	userID := getContextUserID(c)

	listID, err := getValueFromParams(c, "listID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	list, err := h.TodoListService.GetByID(userID, listID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

func (h *Handler) getAllLists(c echo.Context) error {
	userID := getContextUserID(c)

	orderBy := c.QueryParam("sort_by")
	
	lists, err := h.TodoListService.GetAll(userID, &orderBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, resourceResponse{
		Count:      len(lists),
		Results:    lists,
		Pagination: nil,
	})
}

func (h *Handler) createList(c echo.Context) error {
	userID := getContextUserID(c)

	var input model.CreateTodoListDTO
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	id, err := h.TodoListService.Create(userID, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id": id,
	})
}
