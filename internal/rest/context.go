package rest

import (
	"context"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (api *API) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (api *API) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")

	}
	return user
}
