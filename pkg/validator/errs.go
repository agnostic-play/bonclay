package validator

import "fmt"

type BadRequestErrors struct {
	Errors map[string]string `json:"errors"`
}

func (e BadRequestErrors) Error() string {
	return fmt.Sprintf("bad request")
}

func (e BadRequestErrors) GetBadRequestPayload() map[string]string {
	return e.Errors
}
