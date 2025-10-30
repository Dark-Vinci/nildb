package errors

import "errors"

var (
	ErrFileDoesNotExist       = errors.New("file does not exist")
	ErrFileNotOpened          = errors.New("file is not opened")
	ErrFilePathISNil          = errors.New("file path is nil")
	ErrInvalidWhence          = errors.New("invalid whence")
	ErrInvalidPointerPosition = errors.New("invalid pointer position")
)
