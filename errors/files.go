package errors

import perrors "errors"

var (
	ErrFileNotFound = perrors.New("file not found")
)
