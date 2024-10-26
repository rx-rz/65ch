package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rx-rz/65ch/internal/data"
	"io"
	"net/http"
)

type envelope map[string]any

func (api *API) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request input contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("request input contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("request input contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("request input contains incorrect JSON type (at character %d, expected type %T)", unmarshalTypeError.Offset, unmarshalTypeError.Type)
		case errors.Is(err, io.EOF):
			return errors.New("request input must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(dest)
	if err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}
	return nil
}

func (api *API) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) {
	w.Header().Set("Content-Type", "application/json")
	for k, v := range headers {
		w.Header()[k] = v
	}
	js, err := json.Marshal(data)
	if err != nil {
		api.logger.PrintError(err, nil)
	}

	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		api.logger.PrintError(err, nil)
	}
}

func (api *API) handleDBError(w http.ResponseWriter, r *http.Request, err error, message string) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		api.notFoundErrorResponse(w, r)
	case errors.Is(err, data.ErrEditConflict):
		api.conflictResponse(w, r, message)
	default:
		api.internalServerErrorResponse(w, r, err)
	}

}

func (api *API) returnDBError(w http.ResponseWriter, r *http.Request, err error, message string) {
	switch {
	case errors.Is(err, data.ErrRecordNotFound):
		api.notFoundErrorResponse(w, r)
	case errors.Is(err, data.ErrEditConflict):
		api.conflictResponse(w, r, message)
	case errors.Is(err, data.ErrCheckConstraint):
		api.badRequestErrorResponse(w, r, message)
	case errors.Is(err, data.ErrDuplicateKey):
		api.badRequestErrorResponse(w, r, message)
	case errors.Is(err, data.ErrForeignKeyViolation):
		api.badRequestErrorResponse(w, r, message)
	case errors.Is(err, data.ErrInvalidInput):
		api.badRequestErrorResponse(w, r, message)
	default:
		api.internalServerErrorResponse(w, r, err)
	}
}
