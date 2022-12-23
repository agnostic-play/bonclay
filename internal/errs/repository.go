package errs

import "fmt"

var (
	ErrNotFound = fmt.Errorf("data not found")
)
