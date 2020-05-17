package appvalidator

import (
	"github.com/go-playground/validator/v10"
)

func NewValidate() (v *validator.Validate) {
	v = validator.New()
	return v
}
