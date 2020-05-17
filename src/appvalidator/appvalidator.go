package appvalidator

import (
	"github.com/go-playground/validator/v10"
)

// NewValidate creates a validator.Validate with custom validation tags
func NewValidate() (v *validator.Validate) {
	v = validator.New()
	// TODO: add a custom validation tag
	return
}
