package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationErrors formats validator errors into user-friendly messages
func HandleValidationErrors(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationError

		for _, e := range validationErrors {
			field := e.Field()
			message := getErrorMessage(e)

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Handle JSON binding errors
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Invalid request data",
		"error":   err.Error(),
	})
}

// getErrorMessage returns a user-friendly error message based on validation tag
func getErrorMessage(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must not exceed " + e.Param() + " characters"
	case "gt":
		return field + " must be greater than " + e.Param()
	case "gte":
		return field + " must be greater than or equal to " + e.Param()
	case "lt":
		return field + " must be less than " + e.Param()
	case "lte":
		return field + " must be less than or equal to " + e.Param()
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "numeric":
		return field + " must be a number"
	case "alphanum":
		return field + " must contain only alphanumeric characters"
	case "alpha":
		return field + " must contain only alphabetic characters"
	case "oneof":
		return field + " must be one of: " + e.Param()
	case "uuid":
		return field + " must be a valid UUID"
	case "len":
		return field + " must be exactly " + e.Param() + " characters"
	default:
		return field + " is invalid"
	}
}
