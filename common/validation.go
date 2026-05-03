package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns a slice of error messages if any.
func ValidateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag()))
		}
		return errors
	}
	return nil
}

// FormatValidationError joins multiple error messages into a single string.
func FormatValidationError(errors []string) string {
	return strings.Join(errors, ", ")
}
