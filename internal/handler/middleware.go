package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const ctxUserID = "userID"

func (h *Handler) JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")

		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "empty auth header")
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth token")
		}

		if headerParts[1] == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "no token provided")
		}

		claims, err := h.UserService.ParseToken(headerParts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		expires := claims["expires"]
		expiresTime, _ := time.Parse(time.RFC3339, expires.(string))

		if time.Now().UTC().After(expiresTime) {
			return echo.NewHTTPError(http.StatusUnauthorized, "token is expired")
		}

		userID := claims["userID"]

		c.Set(ctxUserID, userID)

		return next(c)
	}
}
