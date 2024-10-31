package rest

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rx-rz/65ch/internal/config"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/jsonlog"
	"net/http"
	"time"
)

type API struct {
	router *httprouter.Router
	models data.Models
	logger *jsonlog.Logger
}

func InitializeAPI(cfg *config.Config) *http.Server {
	api := &API{
		router: httprouter.New(),
		models: data.NewModels(cfg.DB),
		logger: cfg.Logger,
	}

	api.initializeUserRoutes()
	api.initializeCategoryRoutes()
	api.initializeTagRoutes()

	return &http.Server{
		Handler:      api.router,
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
}
