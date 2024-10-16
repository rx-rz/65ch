package rest

import (
	"net/http"
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
