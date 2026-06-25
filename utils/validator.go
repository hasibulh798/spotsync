package utils

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the go-playground validator.
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the structure of the input payload using validator tags.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
