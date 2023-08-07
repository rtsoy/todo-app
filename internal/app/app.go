package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/rtsoy/todo-app/config"
	_ "github.com/rtsoy/todo-app/docs"
	"github.com/rtsoy/todo-app/internal/handler"
	"github.com/rtsoy/todo-app/internal/repository"
	"github.com/rtsoy/todo-app/internal/service"
	"github.com/rtsoy/todo-app/pkg/logger"
	"github.com/rtsoy/todo-app/pkg/postgresql"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Run(cfg *config.Config) {
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log := logger.New(logrus.InfoLevel, e)

	db, err := postgresql.New(cfg)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s", err.Error())
	}

	rpstry := repository.NewRepository(db)
	svc := service.NewService(rpstry)
	hndlr := handler.NewHandler(svc)

	hndlr.InitRoutes(e)

	log.Fatal(e.Start(cfg.HTTPPort))
}
