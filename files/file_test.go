package files

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOperations(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() *File
		action      func(*File) error
		verify      func(*testing.T, *File)
		expectError bool
	}{
		{
			name: "Create and write to a file with success",

			setup: func() *File {
				path := filepath.Join(tempDir, "create_write.txt")
				f := NewFile(path)

				if _, err := f.Create(); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return f
			},

			action: func(f *File) error {
				data := []byte("hello world")
				n, err := f.Write(data)

				if err != nil {
					return err
				}

				if n != len(data) {
					t.Errorf("expected to write %d bytes, wrote %d", len(data), n)
				}

				return nil
			},

			verify: func(t *testing.T, f *File) {
				_, err := f.Seek(0, io.SeekStart)
				if err != nil {
					t.Fatalf("seek failed: %v", err)
				}

				buf := make([]byte, 11)
				n, err := f.Read(buf)

				if err != nil {
					t.Fatalf("read failed: %v", err)
				}

				got := string(buf[:n])
				want := "hello world"

				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			},
		},

		{
			name: "Truncate file content",

			setup: func() *File {
				path := filepath.Join(tempDir, "truncate.txt")
				f := NewFile(path)

				if _, err := f.Create(); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				if _, err := f.Write([]byte("data to truncate")); err != nil {
					t.Fatalf("failed to write: %v", err)
				}

				return f
			},

			action: func(f *File) error {
				return f.Truncate()
			},

			verify: func(t *testing.T, f *File) {
				stat, err := os.Stat(f.path)
				if err != nil {
					t.Fatalf("failed to stat file: %v", err)
				}

				if stat.Size() != 0 {
					t.Errorf("expected file size 0, got %d", stat.Size())
				}
			},
		},

		{
			name: "Sync should not return error",

			setup: func() *File {
				path := filepath.Join(tempDir, "sync.txt")
				f := NewFile(path)
				if _, err := f.Create(); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return f
			},

			action: func(f *File) error {
				return f.Sync()
			},
		},

		{
			name: "Close should close file",

			setup: func() *File {
				path := filepath.Join(tempDir, "close.txt")
				f := NewFile(path)

				if _, err := f.Create(); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				return f
			},

			action: func(f *File) error {
				return f.Close()
			},
		},

		{
			name: "Remove should delete file",

			setup: func() *File {
				path := filepath.Join(tempDir, "remove.txt")
				f := NewFile(path)

				if _, err := f.Create(); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				return f
			},

			action: func(f *File) error {
				return f.Remove()
			},

			verify: func(t *testing.T, f *File) {
				if _, err := os.Stat(f.path); !os.IsNotExist(err) {
					t.Errorf("expected file to be deleted, got err=%v", err)
				}
			},
		},

		{
			name: "Open existing file",

			setup: func() *File {
				path := filepath.Join(tempDir, "open_existing.txt")
				err := os.WriteFile(path, []byte("existing"), 0644)
				if err != nil {
					t.Fatalf("failed to pre-create file: %v", err)
				}
				return NewFile(path)
			},

			action: func(f *File) error {
				_, err := f.Open()
				return err
			},

			verify: func(t *testing.T, f *File) {
				buf := make([]byte, 8)
				n, err := f.Read(buf)

				if err != nil && err != io.EOF {
					t.Fatalf("read failed: %v", err)
				}

				got := string(buf[:n])
				if got != "existing" {
					t.Errorf("expected 'existing', got %q", got)
				}
			},
		},

		{
			name: "Create with empty path should fail",

			setup: func() *File {
				return NewFile("")
			},

			action: func(f *File) error {
				_, err := f.Create()
				return err
			},

			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.setup()

			err := tt.action(f)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.verify != nil {
				tt.verify(t, f)
			}

			if f.f != nil {
				_ = f.Close()
			}
		})
	}
}
