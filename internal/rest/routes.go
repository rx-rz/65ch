package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/v1/movies", "")
}
