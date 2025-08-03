package dto

import (
	"github.com/go-playground/validator/v10"
)

var (
	globalValidator *validator.Validate
)

// GetValidator returns the global validator instance.
// It ensures that the validator is initialized only once and returns the same instance
// for the entire application.
//
// Returns:
//   - *validator.Validate: The global validator instance
func GetValidator() *validator.Validate {
	if globalValidator == nil {
		globalValidator = validator.New(validator.WithRequiredStructEnabled())
	}
	return globalValidator
}
