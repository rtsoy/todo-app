package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
)

// @Summary Delete a list
// @Description Delete a list by its ID
// @Tags Lists
// @Produce json
// @Security ApiKeyAuth
// @Param listID path string true "List ID"
// @Success 204 "No Content"
// @Failure 400 {object} swaggerErrorResponse
// @Failure 500 {object} swaggerErrorResponse
// @Router /api/lists/{listID} [delete]
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

// @Summary Update a list
// @Description Update a list by its ID
// @Tags Lists
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param listID path string true "List ID"
// @Param input body model.UpdateTodoListDTO true "Updated list data"
// @Success 200 "No Content"
// @Failure 400 {object} swaggerErrorResponse
// @Router /api/lists/{listID} [patch]
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

// @Summary Get a list by ID
// @Description Get a list by its ID
// @Tags Lists
// @Produce json
// @Security ApiKeyAuth
// @Param listID path string true "List ID"
// @Success 200 {object} model.TodoList
// @Failure 400 {object} swaggerErrorResponse
// @Failure 404 {object} swaggerErrorResponse
// @Router /api/lists/{listID} [get]
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

// @Summary Get all lists
// @Description Get all lists
// @Tags Lists
// @Produce json
// @Security ApiKeyAuth
// @Param sort_by query string false "Sort lists by"
// @Success 200 {object} resourceResponse
// @Failure 404 {object} swaggerErrorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c echo.Context) error {
	userID := getContextUserID(c)

	orderBy := c.QueryParam("sort_by")

	orderByPtr := &orderBy
	if orderBy == "" {
		orderByPtr = nil
	}

	lists, err := h.TodoListService.GetAll(userID, orderByPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, resourceResponse{
		Count:      len(lists),
		Results:    lists,
		Pagination: nil,
	})
}

// @Summary Create a list
// @Description Create a new list
// @Tags Lists
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body model.CreateTodoListDTO true "New list data"
// @Success 201 {object} createResponse
// @Failure 400 {object} swaggerErrorResponse
// @Router /api/lists [post]
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
