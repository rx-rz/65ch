package main

import (
	_ "github.com/lib/pq"
	"github.com/rx-rz/65ch/internal/config"
	"github.com/rx-rz/65ch/internal/jsonlog"
	"github.com/rx-rz/65ch/internal/rest"
	"log"
	"os"
)

func main() {
	_, err := config.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := jsonlog.New(os.Stdout, "errors.txt")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.InitializeDB()
	cfg := config.New(db, logger)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	api := rest.InitializeAPI(cfg)
	logger.PrintInfo("starting server", map[string]string{
		"addr": api.Addr,
	})
	err = api.ListenAndServe()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()
}
