package url

import "errors"

var (
	ErrNotFound   = errors.New("link not found")
	ErrInvalidURL = errors.New("invalid URL")
)
