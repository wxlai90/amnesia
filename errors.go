package amnesia

import "errors"

var (
	ErrIdNotFound        = errors.New("`id` not found on document")
	ErrIdNotOfStringType = errors.New("`id` must be of string type")
	ErrDocNotExist       = errors.New("doc does not exist")
)
