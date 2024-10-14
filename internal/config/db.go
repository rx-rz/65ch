package config

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

func InitializeDB() (DB *sql.DB, err error) {
	envs, err := LoadEnvVariables()
	db, err := sql.Open("postgres", envs.DbUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(envs.DbMaxOpenConns)
	db.SetMaxIdleConns(envs.DbMaxIdleConns)
	duration, err := time.ParseDuration(envs.DbMaxTimeout)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// connect to the db, throw error after a maximum of 5 secs with no connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
