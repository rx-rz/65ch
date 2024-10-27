package rest

import (
	"errors"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"time"
)

type SuccessInfo struct {
	Status     string      `json:"status"`            // "success" or "error"
	Data       any         `json:"data,omitempty"`    // omitted if null
	Message    string      `json:"message,omitempty"` // human-readable message
	Timestamp  string      `json:"timestamp"`         // ISO8601 format
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code      string            `json:"code"`    // application-specific error code
	Status    string            `json:"status"`  // "success" or "error"
	Message   any               `json:"message"` // user-friendly error message
	Timestamp string            `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"` // additional error context
}

type ErrorCode string

const (
	ErrInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrNotFound          ErrorCode = "NOT_FOUND"
	ErrBadRequest        ErrorCode = "BAD_REQUEST"
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrForbidden         ErrorCode = "FORBIDDEN"
	ErrDuplicateEntry    ErrorCode = "DUPLICATE_ENTRY"
	ErrValidation        ErrorCode = "VALIDATION_ERROR"
	ErrDatabaseOperation ErrorCode = "DATABASE_ERROR"
	ErrInternal          ErrorCode = "INTERNAL_ERROR"
)

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
	var dbErr *data.DBError
	if errors.As(err, &dbErr) {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			api.writeErrorResponse(w, http.StatusNotFound, ErrNotFound, err.Error(), dbErr)
		case errors.Is(err, data.ErrEditConflict):
			api.writeErrorResponse(w, http.StatusConflict, ErrDuplicateEntry, err.Error(), dbErr)
		case errors.Is(err, data.ErrCheckConstraint):
			api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), dbErr)
		case errors.Is(err, data.ErrDuplicateKey):
			api.writeErrorResponse(w, http.StatusConflict, ErrDuplicateEntry, err.Error(), dbErr)
		case errors.Is(err, data.ErrForeignKeyViolation):
			api.writeErrorResponse(w, http.StatusBadRequest, ErrBadRequest, err.Error(), dbErr)
		case errors.Is(err, data.ErrInvalidInput):
			api.writeErrorResponse(w, http.StatusUnprocessableEntity, ErrInvalidInput, err.Error(), dbErr)
		default:
			api.internalServerErrorResponse(w, r, err)
		}
	}
}

func (api *API) writeSuccessResponse(w http.ResponseWriter, status int, data any, message string) {
	response := SuccessInfo{
		Status:    "success",
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	api.writeJSON(w, status, response, nil)
}

func (api *API) writeErrorResponse(w http.ResponseWriter, status int, errorCode ErrorCode, message any, err error) {
	details := make(map[string]string)
	var dbError data.DBError
	if errors.As(err, &dbError) {
		details["operation"] = dbError.Operation
		// TODO: change this to ensure details are logged in dev environments only
		details["detail"] = dbError.Detail
	}
	response := ErrorInfo{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Code:      string(errorCode),
		Message:   message,
		Details:   details,
		Status:    "error",
	}
	api.writeJSON(w, status, response, nil)
}

func (api *API) successResponseWithPagination() {}
