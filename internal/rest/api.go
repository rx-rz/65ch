package rest

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"time"
)

type API struct {
	router *httprouter.Router
	models data.Models
}

func InitializeAPI(db *sql.DB) *http.Server {
	api := &API{
		router: httprouter.New(),
		models: data.NewModels(db),
	}

	api.initializeUserRoutes()

	return &http.Server{
		Handler:      api.router,
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
}
