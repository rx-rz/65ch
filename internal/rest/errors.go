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
	err := api.writeJSON(w, status, env, nil)
	if err != nil {
		api.logError(r, err)
		w.WriteHeader(500)
	}
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

func (api *API) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	api.errorResponse(w, r, http.StatusUnprocessableEntity, errors, false)
}

func (api *API) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.logError(r, err)
	api.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request", true)
}
