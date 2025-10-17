package files

import (
	"bytes"
	"errors"
	"io"

	"github.com/dark-vinci/nildb/interfaces"
)

type MemFile struct {
	buf      *bytes.Buffer
	position int
}

var _ interfaces.IOOperator = (*MemFile)(nil)

func (m *MemFile) Write(p []byte) (n int, err error) {
	return m.buf.Write(p)
}

func (m *MemFile) Read(p []byte) (n int, err error) {
	return bytes.NewReader(m.buf.Bytes()[m.position:]).Read(p)
}

func (m *MemFile) Seek(offset int64, whence int) (int64, error) {
	size, lastPosition := int64(m.buf.Len()), int64(m.position)

	switch whence {
	case io.SeekStart:
		lastPosition = offset
	case io.SeekCurrent:
		lastPosition += offset
	case io.SeekEnd:
		lastPosition = size + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if lastPosition < 0 {
		return 0, errors.New("invalid position")
	}

	m.position = int(lastPosition)

	return lastPosition, nil
}

func (m *MemFile) Close() error {
	if m.buf == nil {
		return errors.New("file already closed")
	}

	m.buf = nil // drop reference so it's unusable

	return nil
}

func (m *MemFile) Remove() error {
	_ = m.Truncate()

	return nil
}

func (m *MemFile) Truncate() error {
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
