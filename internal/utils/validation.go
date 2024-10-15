package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}

func GetValidationErrors(err error) map[string]string {
	errorMessages := make(map[string]string)

	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil {
		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = getErrorMessage(fieldError)
		}
	}
	return errorMessages
}
