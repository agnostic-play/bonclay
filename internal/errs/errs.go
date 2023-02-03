package errs

import "fmt"

type CustomError struct {
	HttpCode int
	Message  string
	Err      error
}

func (err CustomError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func BadRequest(err string) CustomError {
	return CustomError{
		Message:  "bad request",
		HttpCode: 400,
		Err:      fmt.Errorf("%s", err),
	}
}
