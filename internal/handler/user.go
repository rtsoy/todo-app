package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/internal/model"
)

type signInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}

	token, err := h.UserService.GenerateToken(input.Email, input.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (h *Handler) signUp(c echo.Context) error {
	var input model.CreateUserDTO

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}

	id, err := h.UserService.CreateUser(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}
