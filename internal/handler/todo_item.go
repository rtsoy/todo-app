package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
)

// @Summary Delete an item
// @Description Delete an item by its ID
// @Tags Items
// @Produce json
// @Security ApiKeyAuth
// @Param itemID path string true "Item ID"
// @Success 204 "No Content"
// @Failure 400 {object} swaggerErrorResponse
// @Failure 500 {object} swaggerErrorResponse
// @Router /api/lists/{listID}/items/{itemID} [delete]
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

// @Summary Update an item
// @Description Update an item by its ID
// @Tags Items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param itemID path string true "Item ID"
// @Param input body model.UpdateTodoItemDTO true "Updated item data"
// @Success 200 "No Content"
// @Failure 400 {object} swaggerErrorResponse
// @Router /api/lists/{listID}/items/{itemID} [patch]
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

// @Summary Get an item by ID
// @Description Get an item by its ID
// @Tags Items
// @Produce json
// @Security ApiKeyAuth
// @Param itemID path string true "Item ID"
// @Success 200 {object} model.TodoItem
// @Failure 400 {object} swaggerErrorResponse
// @Failure 404 {object} swaggerErrorResponse
// @Router /api/lists/{listID}/items/{itemID} [get]
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

// @Summary Get all items
// @Description Get all items for a specific list
// @Tags Items
// @Produce json
// @Security ApiKeyAuth
// @Param listID path string true "List ID"
// @Param sort_by query string false "Sort items by"
// @Param pagination query model.Pagination false "Pagination options"
// @Success 200 {object} resourceResponse
// @Failure 400 {object} swaggerErrorResponse
// @Failure 404 {object} swaggerErrorResponse
// @Router /api/lists/{listID}/items [get]
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

	orderBy := c.QueryParam("sort_by")

	orderByPtr := &orderBy
	if orderBy == "" {
		orderByPtr = nil
	}

	items, err := h.TodoItemService.GetAll(userID, listID, &pagination, orderByPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, resourceResponse{
		Count:      len(items),
		Results:    items,
		Pagination: &pagination,
	})
}

// @Summary Create an item
// @Description Create a new item for a specific list
// @Tags Items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param listID path string true "List ID"
// @Param input body model.CreateTodoItemDTO true "New item data"
// @Success 201 {object} createResponse
// @Failure 400 {object} swaggerErrorResponse
// @Router /api/lists/{listID}/items [post]
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

	return c.JSON(http.StatusCreated, createResponse{ID: id.String()})
}
