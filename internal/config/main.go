package config

import (
	"database/sql"
)

type Config struct {
	DB *sql.DB
}

func New(db *sql.DB) *Config {
	return &Config{db}
}
