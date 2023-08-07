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

type signInResponse struct {
	Token string `json:"token"`
}

// @Summary Sign In
// @Description Authenticate user using email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body signInInput true "Authentication data"
// @Success 200 {object} signInResponse
// @Failure 400 {object} swaggerErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	token, err := h.UserService.GenerateToken(input.Email, input.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, signInResponse{Token: token})
}

// @Summary Sign Up
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body model.CreateUserDTO true "Registration data"
// @Success 200 {object} createResponse
// @Failure 400 {object} swaggerErrorResponse
// @Failure 409 {object} swaggerErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c echo.Context) error {
	var input model.CreateUserDTO

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	id, err := h.UserService.CreateUser(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusOK, createResponse{ID: id.String()})
}
