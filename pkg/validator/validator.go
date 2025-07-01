package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/id"

	"reflect"
	"strings"
)

func (r *RequestValidator) Validate(i interface{}) error {
	return r.Validator.Struct(i)
}

func (cv *CustomValidator) Validate(i interface{}) (errs error) {
	english := en.New()
	uni := ut.New(english, english)

	cv.translator, _ = uni.GetTranslator("en")
	if err := translations.RegisterDefaultTranslations(cv.Validator, cv.translator); err != nil {
		return
	}

	cv.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	cv.registerCustomValidation()

	if err := cv.Validator.Struct(i); err != nil {
		tmpErrsType := map[string]string{}

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("bad request")
		}

		validatorErrs := err.(validator.ValidationErrors)
		for _, e := range validatorErrs {
			tmpErrsType[strings.ToLower(e.Field())] = fmt.Sprintf(e.Translate(cv.translator))
		}
		if len(tmpErrsType) > 0 {
			return BadRequestErrors{Errors: tmpErrsType}
		}

	}
	return nil
}

func (cv *CustomValidator) addTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = cv.Validator.RegisterTranslation(tag, cv.translator, registerFn, transFn)
}
