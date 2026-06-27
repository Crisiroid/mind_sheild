package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator implements Echo's Validator interface
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

// Validate validates the input struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
