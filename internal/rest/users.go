package rest

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
	"os"
	"time"
)

func (api *API) initializeUserRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/users/login", api.loginUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/users/logout", api.logoutUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/users/register", api.registerUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/users/request-password-reset", api.resetPasswordTokenHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/users/reset-password", api.resetPasswordFormHandler)
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/me", api.authorizedAccessOnly(api.updateUserDetailsHandler))
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/me/email", api.authorizedAccessOnly(api.updateUserEmailHandler))
	api.router.HandlerFunc(http.MethodPatch, "/v1/users/me/password", api.authorizedAccessOnly(api.updateUserPasswordHandler))
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
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, utils.GetValidationErrors(validationError), validationError)
		return
	}
	hashedPassword, _ := utils.HashPassword(req.Password)
	user, _ := api.models.Users.FindByEmail(req.Email)
	if user != nil {
		api.writeErrorResponse(w, http.StatusConflict, ErrDuplicateEntry, "User with provided email already exists", nil)
		return
	}
	user = &data.User{
		Email:         req.Email,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Bio:           req.Bio,
		ProfilePicUrl: req.ProfilePictureUrl,
	}

	err = api.models.Users.Create(user)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusCreated, nil, "User registered successfully")
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

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (api *API) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	err := api.readJSON(w, r, &req)
	if err != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, "", validationError)
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, "Invalid details provided", err)
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
		api.writeErrorResponse(w, http.StatusInternalServerError, ErrInternal, "Token generation failed", err)
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
	api.writeSuccessResponse(w, http.StatusOK, envelope{"token": token}, "Login successful")
}

type UpdateUserDetailsRequest struct {
	FirstName         *string `json:"first_name" validate:"min=1,max=255"`
	LastName          *string `json:"last_name" validate:"min=1,max=255"`
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
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, "", validationError)
		return
	}
	user, err := api.models.Users.FindByID(req.ID)
	if err != nil {
		api.handleDBError(w, r, err)
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
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, nil, "User updated successfully")
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
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, utils.GetValidationErrors(validationError), validationError)
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, "Invalid details provided", nil)
		return
	}
	err = api.models.Users.UpdateEmail(req.Email, req.NewEmail)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, nil, "User email updated successfully")
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
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, utils.GetValidationErrors(validationError), validationError)
		return
	}

	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Email)
	if !matches {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, "Invalid details provided", nil)
		return
	}

	hashedPassword, _ := utils.HashPassword(req.NewPassword)
	err = api.models.Users.UpdatePassword(req.Email, hashedPassword)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, nil, "User password updated successfully")

}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (api *API) resetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, utils.GetValidationErrors(validationError), validationError)
		return
	}
	user, err := api.models.Users.FindByEmail(req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
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
			api.handleDBError(w, r, err)
			return
		}
	} else {
		err = api.models.ResetTokens.Update(&data.ResetToken{
			ResetToken: resetToken,
			Expiration: expiration,
			UserID:     user.ID,
		})
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"reset_token": resetToken}, "User email updated successfully")
}

type ResetPasswordFormRequest struct {
	ResetToken  string `json:"reset_token" validate:"required,reset_token"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (api *API) resetPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordFormRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), err)
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, utils.GetValidationErrors(validationError), validationError)
		return
	}
	existingResetToken, err := api.models.ResetTokens.GetByToken(req.ResetToken)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	user, err := api.models.Users.FindByEmail(existingResetToken.UserID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	err = api.models.Users.UpdatePassword(user.Email, req.NewPassword)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	err = api.models.ResetTokens.Delete(user.ID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, nil, "User password reset successfully")
}
