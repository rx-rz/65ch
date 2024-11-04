package rest

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/rx-rz/65ch/internal/data"
	"github.com/rx-rz/65ch/internal/utils"
	"net/http"
	"os"
	"time"
)

func (api *API) initializeUserRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/login", api.loginUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/logout", api.logoutUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/register", api.registerUserHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/auth/request-password-reset", api.resetPasswordRequestHandler)
	api.router.HandlerFunc(http.MethodPatch, "/v1/auth/reset-password", api.resetPasswordHandler)
	api.router.HandlerFunc(http.MethodGet, "/v1/users/:id", api.authorizedAccessOnly(func(w http.ResponseWriter, r *http.Request) {
		ps := httprouter.ParamsFromContext(r.Context())
		api.getUserDetailsHandler(w, r, ps)
	}))
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
	ctx, cancel := api.CreateContext()
	defer cancel()
	var req CreateUserRequest

	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		api.internalServerErrorResponse(w, r, err)
		return
	}
	user, err := api.models.Users.GetByEmail(ctx, req.Email)
	if user != nil {
		api.conflictResponse(w, "User with email already exists")
		return
	}
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		api.handleDBError(w, r, err)
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

	_, err = api.models.Users.Create(ctx, user)
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
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req LoginUserRequest

	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}

	user, err := api.models.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.badRequestResponse(w, err, "Invalid details provided")
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
	api.writeSuccessResponse(w, http.StatusOK, envelope{"token": token}, "Login successful")
}

type GetUserDetailsRequest struct {
	ID string `json:"id" validate:"required"`
}

func (api *API) getUserDetailsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	req := GetUserDetailsRequest{
		ID: ps.ByName("id"),
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	user, err := api.models.Users.GetByID(ctx, req.ID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	userDetails := map[string]string{
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"email":               user.Email,
		"bio":                 user.Bio,
		"profile_picture_url": user.ProfilePicUrl,
		"created_at":          user.CreatedAt.UTC().String(),
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"user": userDetails}, "")
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
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req UpdateUserDetailsRequest
	v := validator.New()
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	user, err := api.models.Users.GetByID(ctx, req.ID)
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
	updateInfo, err := api.models.Users.UpdateDetails(ctx, user)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "User updated successfully")
}

type UpdateUserEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	NewEmail string `json:"new_email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (api *API) updateUserEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req UpdateUserEmailRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	user, err := api.models.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.Password, user.Password)
	if !matches {
		api.badRequestResponse(w, err, "Invalid details provided")
		return
	}
	updateInfo, err := api.models.Users.UpdateEmail(ctx, req.Email, req.NewEmail)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "User email updated successfully")
}

type UpdateUserPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

func (api *API) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req UpdateUserPasswordRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}

	user, err := api.models.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	matches := utils.CheckPasswordHash(req.CurrentPassword, user.Password)
	if !matches {
		api.badRequestResponse(w, err, "Invalid details provided")
		return
	}

	hashedPassword, _ := utils.HashPassword(req.NewPassword)
	updateInfo, err := api.models.Users.UpdatePassword(ctx, req.Email, hashedPassword)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "User password updated successfully")

}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (api *API) resetPasswordRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req ResetPasswordRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	user, err := api.models.Users.GetByEmail(ctx, req.Email)
	if user == nil {
		api.writeSuccessResponse(w, http.StatusOK, nil, "Token has been sent to your email if you have an account")
		return
	}
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}

	resetToken, expiration := utils.GenerateResetToken()
	previousResetToken, err := api.models.ResetTokens.GetByUserID(ctx, user.ID)

	if previousResetToken == nil {
		_, err = api.models.ResetTokens.Create(ctx, &data.ResetToken{
			ResetToken: resetToken,
			Expiration: expiration,
			UserID:     user.ID,
		})
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
	} else {
		_, err = api.models.ResetTokens.Update(ctx, &data.ResetToken{
			ResetToken: resetToken,
			Expiration: expiration,
			UserID:     user.ID,
		})
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"reset_token": resetToken}, "Token has been sent to your email if you have an account")
}

type ResetPasswordFormRequest struct {
	ResetToken  string `json:"reset_token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (api *API) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req ResetPasswordFormRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	existingResetToken, err := api.models.ResetTokens.GetByToken(ctx, req.ResetToken)

	if existingResetToken.Expiration.UTC().Before(time.Now().UTC()) {
		_, err = api.models.ResetTokens.DeleteByToken(ctx, req.ResetToken)
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
		api.writeErrorResponse(w, http.StatusGone, ErrExpired, "Password reset token has expired", nil)
		return
	}

	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	user, err := api.models.Users.GetByID(ctx, existingResetToken.UserID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	updateInfo, err := api.models.Users.UpdatePassword(ctx, user.Email, hashedPassword)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	_, err = api.models.ResetTokens.DeleteByUserId(ctx, user.ID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "User password reset successfully")
}

func (api *API) followUserHandler(w http.ResponseWriter, r *http.Request) {

}
