package app

import (
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rtsoy/todo-app/config"
	"github.com/rtsoy/todo-app/internal/handler"
	"github.com/rtsoy/todo-app/pkg/postgresql"
)

func Run(cfg *config.Config) {
	e := echo.New()

	db, err := postgresql.New(cfg)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %v", err)
	}
	_ = db

	hndlr := handler.NewHandler()
	hndlr.InitRoutes(e)

	e.Logger.Fatal(e.Start(cfg.HTTPPort))
}
