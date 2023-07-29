package main

import (
	"log"

	"github.com/rtsoy/todo-app/config"
	"github.com/rtsoy/todo-app/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
