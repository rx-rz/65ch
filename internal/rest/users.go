package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
)

func (api *API) initializeUserRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/users", api.registerUserHandler)
}

type CreateUser struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1"`
	LastName  string `json:"last_name" validate:"required,min=1"`
}

func (api *API) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var user CreateUser
	v := validator.New()
	err := api.readJSON(w, r, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = v.Struct(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	hashedPassword, _ := utils.HashPassword(user.Password)
	dbUser := data.User{
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	err = api.models.Users.Create(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = api.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		http.Error(w, "error occured", http.StatusInternalServerError)
		return
	}

}
