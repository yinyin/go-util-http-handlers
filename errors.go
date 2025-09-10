package httphandlers

import (
	"errors"
)

// ErrEmptyRequestBody indicate request body is empty.
var ErrEmptyRequestBody = errors.New("empty request body")

// ErrRequireContentType indicate Content-Type header is missing.
var ErrRequireContentType = errors.New("require Content-Type header")
