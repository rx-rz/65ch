package main

import (
	_ "github.com/lib/pq"
	"github.com/rx-rz/65ch/internal/config"
	"github.com/rx-rz/65ch/internal/rest"
	"log"
)

func main() {
	_, err := config.LoadEnvVariables()

	if err != nil {
		log.Fatal(err)
	}
	db, err := config.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	api := rest.InitializeAPI(db)
	err = api.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing the database: %s", err)
		}
	}()
}
