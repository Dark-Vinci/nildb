package files

import (
	"bytes"
	"io"

	"github.com/dark-vinci/nildb/errors"
	"github.com/dark-vinci/nildb/interfaces"
)

type MemFile struct {
	buf      *bytes.Buffer
	position int
}

// majorly for test and simulate in memory database
var _ interfaces.IOOperator = (*MemFile)(nil)

func (m *MemFile) Write(p []byte) (n int, err error) {
	if m.buf == nil {
		m.buf = new(bytes.Buffer)
	}

	return m.buf.Write(p)
}

func (m *MemFile) Read(p []byte) (n int, err error) {
	if m.buf == nil {
		m.buf = new(bytes.Buffer)
	}

	return bytes.NewReader(m.buf.Bytes()[m.position:]).Read(p)
}

func (m *MemFile) Seek(offset int64, whence int) (int64, error) {
	if m.buf == nil {
		m.buf = new(bytes.Buffer)
	}

	size, lastPosition := int64(m.buf.Len()), int64(m.position)

	switch whence {
	case io.SeekStart:
		lastPosition = offset
	case io.SeekCurrent:
		lastPosition += offset
	case io.SeekEnd:
		lastPosition = size + offset
	default:
		return 0, errors.ErrInvalidWhence
	}

	if lastPosition < 0 {
		return 0, errors.ErrInvalidPointerPosition
	}

	m.position = int(lastPosition)

	return lastPosition, nil
}

func (m *MemFile) Close() error {
	if m.buf != nil {
		m.buf = nil // drop reference so it's unusable
	}

	return nil
}

func (m *MemFile) Remove() error {
	_ = m.Truncate()

	return nil
}

func (m *MemFile) Truncate() error {
	if m.buf != nil {
		m.buf = new(bytes.Buffer)
	}

	m.buf.Reset()
	m.position = 0

	return nil
}

func (m *MemFile) Sync() error {
	return nil
}

func (m *MemFile) Create() (interfaces.IOOperator, error) {
	return &MemFile{
		buf:      &bytes.Buffer{},
		position: 0,
	}, nil
}

func (m *MemFile) Open() (interfaces.IOOperator, error) {
	return &MemFile{
		buf:      &bytes.Buffer{},
		position: 0,
	}, nil
}
