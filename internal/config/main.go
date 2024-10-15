package config

import (
	"database/sql"
	"github.com/rx-rz/65ch/internal/jsonlog"
)

type Config struct {
	DB     *sql.DB
	Logger *jsonlog.Logger
}

func New(db *sql.DB, logger *jsonlog.Logger) *Config {
	return &Config{db, logger}
}
