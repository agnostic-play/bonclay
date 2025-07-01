package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

type CustomValidator struct {
	Validator  *validator.Validate
	translator ut.Translator
}

type RequestValidator struct {
	Validator *validator.Validate
}
