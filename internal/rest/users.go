package rest

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
	"os"
	"time"
)

func (api *API) initializeUserRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/register", api.registerUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/login", api.loginUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/logout", api.logoutUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/reset-password/request", api.resetPasswordTokenHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/reset-password/reset", api.resetPasswordFormHandler)
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/update", api.authorizedAccessOnly(api.updateUserDetailsHandler))
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/update-email", api.authorizedAccessOnly(api.updateUserEmailHandler))
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/update-password", api.authorizedAccessOnly(api.updateUserPasswordHandler))
}

type CreateUserRequest struct {
	Email             string `json:"email" validate:"required,email,max=255"`
	Password          string `json:"password" validate:"required,min=8,max=72"`
	FirstName         string `json:"first_name" validate:"required,min=1,max=255"`
	LastName          string `json:"last_name" validate:"required,min=1,max=255"`
	Bio               string `json:"bio"`
	ProfilePictureUrl string `json:"profile_picture_url"`
}

func (api *API) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	v := validator.New()
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	validationError := v.Struct(req)
	if validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	hashedPassword, _ := utils.HashPassword(req.Password)

	user := &data.User{
		Email:         req.Email,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Bio:           req.Bio,
		ProfilePicUrl: req.ProfilePictureUrl,
	}

	err = api.models.Users.Create(user)
	if err != nil {
		api.handleDBError(w, r, err, "User with provided email already exists")
		return
	}
	api.writeJSON(w, http.StatusCreated, envelope{"status": "success", "data": nil, "message": "User registered successfully"}, nil)
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (api *API) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (api *API) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		api.handleDBError(w, r, err, "")
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.badRequestErrorResponse(w, r, "Invalid details provided")
		return
	}
	token, err := utils.GenerateToken(map[string]string{
		"email":               user.Email,
		"id":                  user.ID,
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"profile_picture_url": user.ProfilePicUrl,
		"bio":                 user.Bio,
	}, os.Getenv("JWT_SECRET"))
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    fmt.Sprintf("Bearer %s", token),
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": map[string]string{"token": token}}, nil)
}

type UpdateUserDetailsRequest struct {
	FirstName         *string `json:"first_name", validate:"min=1,max=255"`
	LastName          *string `json:"last_name", validate:"min=1,max=255"`
	Bio               *string `json:"bio"`
	ProfilePictureUrl *string `json:"profile_picture_url"`
	Activated         *bool   `json:"activated"`
	ID                string  `json:"id" validate:"required"`
}

func (api *API) updateUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserDetailsRequest
	v := validator.New()
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	user, err := api.models.Users.FindByID(req.ID)
	if err != nil {
		api.handleDBError(w, r, err, "")
		return
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.ProfilePictureUrl != nil {
		user.ProfilePicUrl = *req.ProfilePictureUrl
	}
	if req.Activated != nil {
		user.Activated = *req.Activated
	}
	_, err = api.models.Users.UpdateDetails(user)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	api.writeJSON(w, http.StatusOK, envelope{"data": map[string]any{
		"user": map[string]string{
			"email": user.Email,
			"id":    user.ID,
		},
	}, "status": "success"}, nil)

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
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			api.notFoundErrorResponse(w, r)
			return
		}
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.badRequestErrorResponse(w, r, "Invalid details provided")
		return
	}
	err = api.models.Users.UpdateEmail(req.Email, req.NewEmail)
	if err != nil {
		api.handleDBError(w, r, err, "")
		return
	}
	api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": nil, "message": "User email updated successfully"}, nil)
}

type UpdateUserPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"password"`
	NewPassword string `json:"new_password" validate:"new_password"`
}

func (api *API) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserPasswordRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}

	user, err := api.models.Users.FindByEmail(req.Email)
	matches := utils.CheckPasswordHash(req.Password, user.Email)
	if !matches {
		api.badRequestErrorResponse(w, r, "invalid details provided")
	}
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			api.notFoundErrorResponse(w, r)
			return
		}
	}
	hashedPassword, _ := utils.HashPassword(req.NewPassword)
	err = api.models.Users.UpdatePassword(req.Email, hashedPassword)
	if err != nil {
		api.handleDBError(w, r, err, "")
	}
	api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": nil, "message": "User email updated successfully"}, nil)

}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (api *API) resetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) && user == nil {
			api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": nil, "message": "Token sent to user email if user exists"}, nil)
		}
		api.internalServerErrorResponse(w, r, err)
		return
	}
	resetToken, expiration := utils.GenerateResetToken()
	previousResetToken, err := api.models.ResetTokens.GetByUserID(user.ID)

	if previousResetToken == nil {
		err = api.models.ResetTokens.Create(&data.ResetToken{
			ResetToken: resetToken,
			Expiration: expiration,
			UserID:     user.ID,
		})
		if err != nil {
			api.badRequestErrorResponse(w, r, err.Error())
			return
		}
	} else {
		err = api.models.ResetTokens.Update(&data.ResetToken{
			ResetToken: resetToken,
			Expiration: expiration,
			UserID:     user.ID,
		})
		if err != nil {
			api.badRequestErrorResponse(w, r, err.Error())
			return
		}
	}
	api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": resetToken, "message": "Token sent to user email if user exists"}, nil)
}

type ResetPasswordFormRequest struct {
	ResetToken  string `json:"reset_token" validate:"required,reset_token"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (api *API) resetPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordFormRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestErrorResponse(w, r, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, r, utils.GetValidationErrors(validationError))
		return
	}
	existingResetToken, err := api.models.ResetTokens.GetByToken(req.ResetToken)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	user, err := api.models.Users.FindByEmail(existingResetToken.UserID)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	err = api.models.Users.UpdatePassword(user.Email, req.NewPassword)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	err = api.models.ResetTokens.Delete(user.ID)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	api.writeJSON(w, http.StatusOK, envelope{"status": "success", "data": nil, "message": "User password reset successfully"}, nil)
}
