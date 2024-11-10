package rest

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
	"os"
	"strings"
)

type UserClaims struct {
	jwt.RegisteredClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (api *API) authorizedAccessOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.unauthorizedResponse(w, r)
			return
		}
		headerParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(headerParts) == 0 {
			api.unauthorizedResponse(w, r)
			return
		}
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			api.invalidTokenResponse(w, r)
			return
		}
		userClaims := &UserClaims{}
		claims, err := utils.DecodeToken(headerParts[1], os.Getenv("JWT_SECRET"), userClaims)
		if err != nil {
			api.invalidTokenResponse(w, r)
			return
		}
		ctx, cancel := api.CreateContext()
		defer cancel()
		user, err := api.models.Users.GetByID(ctx, claims.ID)
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
		r = api.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	}

}
