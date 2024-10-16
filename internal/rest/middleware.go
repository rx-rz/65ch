package rest

import (
	"net/http"
	"strings"
)

func (api *API) authorizedAccessOnly(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

		}
		next.ServeHTTP(w, r)
	}

}
