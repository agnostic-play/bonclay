package validator

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	customText = "{0} hanya dapat berisi huruf dan angka saja"
	customJSON = "{0} harus berformat json"
	urlPath    = "{0} haya dapat di isi dengan format url"
)

func (cv *CustomValidator) registerCustomValidation() {
	_ = cv.Validator.RegisterValidation("urlPath", ValidateURLPath)
	_ = cv.Validator.RegisterValidation("customText", ValidateText)
	_ = cv.Validator.RegisterValidation("customJSON", ValidateJSON)
	cv.addTranslation("urlPath", fmt.Sprintf(urlPath))
	cv.addTranslation("customText", fmt.Sprintf(customText))
	cv.addTranslation("customJSON", fmt.Sprintf(customJSON))
}

func ValidateJSON(fl validator.FieldLevel) bool {
	return json.Valid([]byte(fl.Field().String()))
}

func ValidateText(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^[\w\s-.,)(]+$`)
	return regex.MatchString(fl.Field().String())
}

func ValidateURLPath(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile("^\\/[/.a-zA-Z0-9-]+$")
	return regex.MatchString(fl.Field().String())
}
