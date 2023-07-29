package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rtsoy/todo-app/config"
)

func New(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PSQLHost, cfg.PSQLPort, cfg.PSQLUser, cfg.PSQLPassword,
		cfg.PSQLDBName, cfg.PSQLSSLMode),
	)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
