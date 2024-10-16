package rest

import (
	"net/http"
	"strings"
)

func (api *API) authorizedAccessOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(headerParts) == 0 {
			api.unauthorizedResponse(w, r)
			return
		}
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			api.invalidTokenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}

}
