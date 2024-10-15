package rest

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
	"os"
)

func (api *API) initializeUserRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/users", api.registerUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/users/login", api.loginUserHandler)
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
	validationError := v.Struct(user)
	if validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
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
		if errors.Is(err, data.ErrEditConflict) {
			api.conflictResponse(w, r, fmt.Sprintf("User with email %s already exists", user.Email))
			return
		} else {
			api.internalServerErrorResponse(w, r, err)
			return
		}

	}
	err = api.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}

}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (api *API) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user LoginUser

	err := api.readJSON(w, r, &user)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
	}

	v := validator.New()
	err = v.Struct(user)
	if err != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(err))
		return
	}
	dbUser, err := api.models.Users.FindByEmail(user.Email)
	matches := utils.CheckPasswordHash(user.Password, dbUser.Password)
	if !matches {
		api.badRequestErrorResponse(w, r, "Invalid details provided")
	}
	token := utils.GenerateToken(map[string]string{
		"email": user.Email,
	}, os.Getenv("JWT_SECRET"))
	r.Header.Add("Authorization", "Bearer "+token)
	err = api.writeJSON(w, http.StatusOK, envelope{"user": map[string]string{
		"email":     user.Email,
		"id":        dbUser.ID,
		"firstName": dbUser.FirstName,
		"lastName":  dbUser.LastName,
	}}, nil)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
	}
}
