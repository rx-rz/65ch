package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rx-rz/65ch/internal/config"
	"log"
	"net/http"
)

func main() {
	_, err := config.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Addr: ":8080",
	}
	db, err := config.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
