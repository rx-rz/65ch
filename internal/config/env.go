package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Env struct {
	DbUrl          string
	Port           string
	Env            string
	DbMaxOpenConns int
	DbMaxIdleConns int
	DbMaxTimeout   string
	JwtSecret      string
}

const (
	DefaultEnv          = "development"
	DefaultPort         = "8080"
	DefaultMaxTimeout   = "30s"
	DefaultMaxOpenConns = 10
	DefaultMaxIdleConns = 5
)

func LoadEnvVariables() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, err
	}

	e := Env{
		DbUrl:          getEnv("DB_URL", ""),
		Port:           getEnv("PORT", DefaultPort),
		JwtSecret:      getEnv("JWT_SECRET", ""),
		DbMaxTimeout:   getEnv("DB_MAX_TIMEOUT", DefaultMaxTimeout),
		Env:            getEnv("ENV", DefaultEnv),
		DbMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", DefaultMaxOpenConns),
		DbMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", DefaultMaxIdleConns),
	}
	return e, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Error converting %s to integer, using default value: %d", key, fallback)
			return fallback
		}
		return result
	}
	return fallback
}
