package app

import (
	"github.com/labstack/echo/v4"
	"github.com/rtsoy/todo-app/config"
	"github.com/rtsoy/todo-app/internal/handler"
)

func Run(cfg *config.Config) {
	e := echo.New()

	hndlr := handler.NewHandler()
	hndlr.InitRoutes(e)

	e.Logger.Fatal(e.Start(cfg.HTTPPort))
}
