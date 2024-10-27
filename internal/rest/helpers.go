package rest

import (
	"encoding/json"
	"errors"
	"fmt"
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
			err = fmt.Errorf("request input contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			err = errors.New("request input contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				err = fmt.Errorf("request input contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			err = fmt.Errorf("request input contains incorrect JSON type (at character %d, expected type %T)", unmarshalTypeError.Offset, unmarshalTypeError.Type)
		case errors.Is(err, io.EOF):
			err = errors.New("request input must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			panic(err)
		}
		return err
	}
	if dec.More() {
		err = errors.New("body must contain a single JSON value")
		return err
	}
	return nil
}

func (api *API) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) {
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
