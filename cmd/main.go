package main

import (
	"log"

	"github.com/rtsoy/todo-app/config"
	"github.com/rtsoy/todo-app/internal/app"
)

// @title       TodoApp API
// @version     1.0
// @description A server for TodoApp Application

// @host localhost:3000
// @Basepath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
