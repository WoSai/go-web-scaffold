package validator

import (
	"github.com/go-playground/validator/v10"
)

var (
	defaultValidator *validator.Validate
)

func Struct(s interface{}) error {
	return defaultValidator.Struct(s)
}

func init() {
	defaultValidator = validator.New()
}
