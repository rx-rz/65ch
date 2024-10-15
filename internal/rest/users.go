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
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/update", api.updateUserDetailsHandler)
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1"`
	LastName  string `json:"last_name" validate:"required,min=1"`
}

func (api *API) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var user CreateUserRequest
	v := validator.New()
	err := api.readJSON(w, r, &user)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
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

		} else {
			api.internalServerErrorResponse(w, r, err)

		}

	}
	err = api.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)

	}

}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (api *API) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user LoginUserRequest

	err := api.readJSON(w, r, &user)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
	}

	v := validator.New()
	err = v.Struct(user)
	if err != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(err))
	}
	dbUser, err := api.models.Users.FindByEmail(user.Email)
	matches := utils.CheckPasswordHash(user.Password, dbUser.Password)
	if !matches {
		api.badRequestErrorResponse(w, r, "Invalid details provided")
	}
	token, err := utils.GenerateToken(map[string]string{
		"email": user.Email,
	}, os.Getenv("JWT_SECRET"))
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	err = api.writeJSON(w, http.StatusOK, envelope{"token": token}, nil)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
	}
}

type UpdateUserDetailsRequest struct {
	FirstName         *string `json:"first_name"`
	LastName          *string `json:"last_name"`
	Bio               *string `json:"bio"`
	ProfilePictureUrl *string `json:"profile_picture_url"`
	Activated         *bool   `json:"activated"`
	ID                string  `json:"id" validate:"required"`
}

func (api *API) updateUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var user UpdateUserDetailsRequest
	v := validator.New()
	err := api.readJSON(w, r, &user)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
	}
	validationError := v.Struct(user)
	if validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
	}
	dbUser, err := api.models.Users.FindByID(user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			api.notFoundErrorResponse(w, r)
		}
		api.internalServerErrorResponse(w, r, err)
	}
	if user.FirstName != nil {
		dbUser.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		dbUser.LastName = *user.LastName
	}
	if user.Bio != nil {
		dbUser.Bio = *user.Bio
	}
	if user.ProfilePictureUrl != nil {
		dbUser.ProfilePicUrl = *user.ProfilePictureUrl
	}
	if user.Activated != nil {
		dbUser.Activated = *user.Activated
	}
	updatedUser, err := api.models.Users.UpdateDetails(dbUser)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
	}
	err = api.writeJSON(w, http.StatusOK, envelope{"data": map[string]any{
		"user":   updatedUser,
		"status": "success",
	}}, nil)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)

	}
}

type UpdateUserEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	NewEmail string `json:"new_email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (api *API) updateUserEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserEmailRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
	}
	v := validator.New()
	validationError := v.Struct(req)
	if validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))

	}
	_, err = api.models.Users.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			api.notFoundErrorResponse(w, r)
		}
	}
	err = api.models.Users.UpdateEmail(req.Email, req.NewEmail)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			api.conflictResponse(w, r, "A user with provided email already exists")
		}
		api.internalServerErrorResponse(w, r, err)
	}
	err = api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": map[string]string{"message": "User email updated successfully"}}, nil)
}
