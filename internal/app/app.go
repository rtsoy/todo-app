package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

const gracefulShutdownTimeout = 3 * time.Second

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

	go func() {
		if err := e.Start(cfg.HTTPPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while starting the echo server: %s", err.Error())
		} else {
			log.Println("Todo App started...")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error while shutting down the server: %s", err.Error())
	} else {
		log.Println("Server shut down gracefully.")
	}

	if err := db.Close(); err != nil {
		log.Fatalf("Error while closing the database: %s", err.Error())
	} else {
		log.Println("Database connection closed.")
	}

	log.Println("Exiting... Have a nice day!")
}
