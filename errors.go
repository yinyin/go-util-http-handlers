package httphandlers

import (
	"errors"
)

// ErrEmptyRequestBody indicate request body is empty.
var ErrEmptyRequestBody = errors.New("empty request body")
