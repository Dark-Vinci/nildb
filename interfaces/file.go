package interfaces

import (
	"io"
)

type IOOperator interface {
	io.Writer
	io.Reader
	io.Seeker
	io.Closer

	Remove() error
	Truncate() error
	Sync() error
	Create() (IOOperator, error)
	Open() (IOOperator, error)
}
