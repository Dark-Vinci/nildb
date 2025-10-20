package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dark-vinci/nildb/interfaces"
)

type File struct {
	path string
	f    *os.File
}

var _ interfaces.IOOperator = (*File)(nil)

func NewFile(path string) *File {
	return &File{
		path: path,
		f:    nil,
	}
}

func (f *File) Write(p []byte) (n int, err error) {
	if f.f == nil {
		return 0, errors.New("file is not open")
	}

	write, err := f.f.Write(p)
	if err != nil {
		fmt.Println("Error: file cannot be written", err)
		return 0, err
	}

	return write, nil
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.f == nil {
		return 0, errors.New("file is not open")
	}

	val, err := f.f.Read(p)
	if err != nil {
		fmt.Println("Error: file cannot be read", err)
		return 0, err
	}

	return val, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.f == nil {
		return 0, errors.New("file is not open")
	}

	n, err := f.f.Seek(offset, whence)
	if err != nil {
		fmt.Println("Error: file cannot be seeked", err)
		return 0, err
	}

	return n, nil
}

func (f *File) Close() error {
	if f.f == nil {
		return errors.New("file is not open")
	}

	if err := f.f.Close(); err != nil {
		fmt.Println("Error: file cannot be closed", err)
		return err
	}

	return nil
}

func (f *File) Remove() error {
	if f.path != "" {
		return errors.New("file path cannot be empty")
	}

	if err := os.Remove(f.path); err != nil {
		fmt.Println("Error: file cannot be removed", err)
		return err
	}

	return nil
}

func (f *File) Truncate() error {
	if f.f == nil {
		return errors.New("file is not open")
	}

	if err := os.Truncate(f.path, 0); err != nil {
		fmt.Println("Error: file cannot be truncated", err)
		return err
	}

	return nil
}

func (f *File) Sync() error {
	if err := f.f.Sync(); err != nil {
		fmt.Println("Error: file cannot be synced", err)
		return err
	}

	return nil
}

func (f *File) Create() (interfaces.IOOperator, error) {
	if f.path == "" {
		return interfaces.IOOperator(nil), errors.New("file path cannot be empty")
	}

	if parent := filepath.Dir(f.path); parent != "" {
		if err := os.MkdirAll(parent, 0755); err != nil {
			fmt.Println("Error: directory cannot be created", err)
			return interfaces.IOOperator(nil), err
		}
	}

	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error: file cannot be created", err)
		return interfaces.IOOperator(nil), err
	}

	f.f = file

	return f, nil
}

func (f *File) Open() (interfaces.IOOperator, error) {
	if f.path == "" {
		return nil, errors.New("file path cannot be empty")
	}

	if f.f != nil {
		return f, nil
	}

	file, err := os.OpenFile(f.path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error: file cannot be opened", err)
		return interfaces.IOOperator(nil), err
	}

	f.f = file

	return f, nil
}
