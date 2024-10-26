package rest

import (
	"errors"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"time"
)

type SuccessInfo struct {
	Status    string      `json:"status"`            // "success" or "error"
	Data      interface{} `json:"data,omitempty"`    // omitted if null
	Message   string      `json:"message,omitempty"` // human-readable message
	Error     *ErrorInfo  `json:"error,omitempty"`   // omitted on success
	RequestID string      `json:"requestId"`         // for tracking/debugging
	Timestamp time.Time   `json:"timestamp"`         // ISO8601 format
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code    string            `json:"code"`              // application-specific error code
	Message string            `json:"message"`           // user-friendly error message
	Details map[string]string `json:"details,omitempty"` // additional error context
}

func (api *API) logError(r *http.Request, err error) {
	api.logger.PrintFatal(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (api *API) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, isServerError bool) {
	// https: //github.com/omniti-labs/jsend
	var env envelope
	if isServerError == true {
		env = envelope{"status": "error", "message": message}
	}
	env = envelope{"status": "fail", "data": message}
	api.writeJSON(w, status, env, nil)
	return
}

func (api *API) notFoundErrorResponse(w http.ResponseWriter, r *http.Request) {
	api.errorResponse(w, r, http.StatusNotFound, "The requested resource could not be found", false)
}

func (api *API) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	api.errorResponse(w, r, http.StatusBadRequest, message, false)
}

func (api *API) conflictResponse(w http.ResponseWriter, r *http.Request, message any) {
	api.errorResponse(w, r, http.StatusConflict, message, true)
}

func (api *API) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors []string) {
	api.errorResponse(w, r, http.StatusUnprocessableEntity, envelope{"errors": errors}, false)
}

func (api *API) invalidTokenResponse(w http.ResponseWriter, r *http.Request) {
	api.errorResponse(w, r, http.StatusUnprocessableEntity, "Invalid token provided", false)
}
func (api *API) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	api.errorResponse(w, r, http.StatusUnauthorized, "You are not authorised to use this resourece", false)
}

func (api *API) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.logError(r, err)
	api.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request", true)
}

func (api *API) handleDBError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, data.ErrRecordNotFound):
		api.notFoundErrorResponse(w, r)
	case errors.Is(err, data.ErrEditConflict):
		api.conflictResponse(w, r, err.Error())
	case errors.Is(err, data.ErrCheckConstraint):
		api.badRequestErrorResponse(w, r, err.Error())
	case errors.Is(err, data.ErrDuplicateKey):
		api.badRequestErrorResponse(w, r, err.Error())
	case errors.Is(err, data.ErrForeignKeyViolation):
		api.badRequestErrorResponse(w, r, err.Error())
	case errors.Is(err, data.ErrInvalidInput):
		api.badRequestErrorResponse(w, r, err.Error())
	default:
		api.internalServerErrorResponse(w, r, err)
	}
}

func (api *API) writeError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	var dbErr *data.DBError
	if errors.As(err, &dbErr) {
		api.logger.PrintError(err, map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
			"error":          dbErr.Error(),
		})
	} else {
		api.logger.PrintError(err, map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
			"error":          err.Error(),
		})
	}
	if status != http.StatusInternalServerError {
		api.writeJSON(w, status, envelope{"message": message, "status": "error"}, nil)
		return
	}
	api.writeJSON(w, status, envelope{"message": message, "status": "fail"}, nil)
	return
}

func (api *API) writeSuccessResponse(w http.ResponseWriter, r *http.Request, data any, message string) {

}
