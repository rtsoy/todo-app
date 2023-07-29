package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" env-default:":8080"`

	PSQLHost     string `env:"POSTGRESQL_HOST"`
	PSQLPort     string `env:"POSTGRESQL_PORT"`
	PSQLUser     string `env:"POSTGRESQL_USER"`
	PSQLPassword string `env:"POSTGRESQL_PASSWORD"`
	PSQLDBName   string `env:"POSTGRESQL_DBNAME"`
	PSQLSSLMode  string `env:"POSTGRESQL_SSL_MODE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	return cfg, err
}
