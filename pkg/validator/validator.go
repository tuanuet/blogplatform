package validator

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Init initializes custom validators
func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators here
		_ = v // Use v to register custom validators
	}
}

// FormatValidationErrors formats validator.ValidationErrors into a user-friendly format
func FormatValidationErrors(err error) ValidationErrors {
	var errors ValidationErrors

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: getValidationMessage(e),
			})
		}
	}

	return errors
}

// getValidationMessage returns a user-friendly error message for a validation error
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters"
	case "max":
		return "Must be at most " + e.Param() + " characters"
	case "uuid":
		return "Must be a valid UUID"
	case "oneof":
		return "Must be one of: " + e.Param()
	case "gt":
		return "Must be greater than " + e.Param()
	case "gte":
		return "Must be greater than or equal to " + e.Param()
	case "lt":
		return "Must be less than " + e.Param()
	case "lte":
		return "Must be less than or equal to " + e.Param()
	default:
		return "Invalid value"
	}
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
