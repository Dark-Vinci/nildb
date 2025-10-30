package files

import (
	"bytes"
	"io"
	"testing"
)

func TestMemFileOperations(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *MemFile
		action      func(*testing.T, *MemFile) error
		verify      func(*testing.T, *MemFile)
		expectError bool
	}{
		{
			name: "Create and write to MemFile successfully",

			setup: func() *MemFile {
				m, err := (&MemFile{}).Create()
				if err != nil {
					t.Fatalf("failed to create MemFile: %v", err)
				}

				return m.(*MemFile)
			},

			action: func(t *testing.T, m *MemFile) error {
				data := []byte("in-memory file test")
				n, err := m.Write(data)
				if err != nil {
					return err
				}

				if n != len(data) {
					t.Errorf("expected to write %d bytes, got %d", len(data), n)
				}

				return nil
			},

			verify: func(t *testing.T, m *MemFile) {
				buf := make([]byte, 64)
				_, err := m.Seek(0, io.SeekStart)
				if err != nil {
					t.Fatalf("seek failed: %v", err)
				}

				n, err := m.Read(buf)
				if err != nil && err != io.EOF {
					t.Fatalf("read failed: %v", err)
				}

				got := string(buf[:n])
				want := "in-memory file test"
				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			},
		},

		{
			name: "Seek and read correctly from middle",

			setup: func() *MemFile {
				m := &MemFile{buf: bytes.NewBufferString("abcdefghij")}
				return m
			},

			action: func(t *testing.T, m *MemFile) error {
				pos, err := m.Seek(5, io.SeekStart)
				if err != nil {
					return err
				}

				if pos != 5 {
					t.Errorf("expected position 5, got %d", pos)
				}

				return nil
			},

			verify: func(t *testing.T, m *MemFile) {
				buf := make([]byte, 5)
				n, err := m.Read(buf)
				if err != nil && err != io.EOF {
					t.Fatalf("read failed: %v", err)
				}

				got := string(buf[:n])
				want := "fghij"

				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			},
		},

		{
			name: "Seek from end",

			setup: func() *MemFile {
				m := &MemFile{buf: bytes.NewBufferString("12345")}
				return m
			},

			action: func(t *testing.T, m *MemFile) error {
				pos, err := m.Seek(-2, io.SeekEnd)
				if err != nil {
					return err
				}

				if pos != 3 {
					t.Errorf("expected position 3, got %d", pos)
				}

				return nil
			},

			verify: func(t *testing.T, m *MemFile) {
				buf := make([]byte, 2)
				n, err := m.Read(buf)

				if err != nil && err != io.EOF {
					t.Fatalf("read failed: %v", err)
				}

				got := string(buf[:n])
				if got != "45" {
					t.Errorf("expected %q, got %q", "45", got)
				}
			},
		},

		{
			name: "Seek with invalid whence should fail",

			setup: func() *MemFile {
				return &MemFile{buf: bytes.NewBufferString("data")}
			},

			action: func(t *testing.T, m *MemFile) error {
				_, err := m.Seek(0, 99) // invalid whence
				return err
			},

			expectError: true,
		},

		{
			name: "Seek to negative position should fail",

			setup: func() *MemFile {
				return &MemFile{buf: bytes.NewBufferString("data")}
			},

			action: func(t *testing.T, m *MemFile) error {
				_, err := m.Seek(-10, io.SeekStart)
				return err
			},

			expectError: true,
		},

		{
			name: "Truncate clears buffer",

			setup: func() *MemFile {
				return &MemFile{buf: bytes.NewBufferString("truncate me")}
			},

			action: func(t *testing.T, m *MemFile) error {
				return m.Truncate()
			},

			verify: func(t *testing.T, m *MemFile) {
				if m.buf.Len() != 0 {
					t.Errorf("expected empty buffer after truncate, got %d", m.buf.Len())
				}
			},
		},

		{
			name: "Close clears buffer reference",

			setup: func() *MemFile {
				return &MemFile{buf: bytes.NewBufferString("close test")}
			},

			action: func(t *testing.T, m *MemFile) error {
				return m.Close()
			},

			verify: func(t *testing.T, m *MemFile) {
				if m.buf != nil {
					t.Errorf("expected buffer to be nil after close")
				}
			},
		},

		{
			name: "Remove resets buffer",

			setup: func() *MemFile {
				return &MemFile{buf: bytes.NewBufferString("remove me")}
			},

			action: func(t *testing.T, m *MemFile) error {
				return m.Remove()
			},

			verify: func(t *testing.T, m *MemFile) {
				if m.buf.Len() != 0 {
					t.Errorf("expected buffer to be empty after remove")
				}
			},
		},

		{
			name: "Sync returns nil",

			setup: func() *MemFile {
				return &MemFile{}
			},

			action: func(t *testing.T, m *MemFile) error {
				return m.Sync()
			},
		},

		{
			name: "Open returns a new MemFile",

			setup: func() *MemFile {
				return &MemFile{}
			},

			action: func(t *testing.T, m *MemFile) error {
				ioOp, err := m.Open()
				if err != nil {
					return err
				}

				newM := ioOp.(*MemFile)
				if newM.buf == nil {
					t.Errorf("expected initialized buffer after open")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()

			err := tt.action(t, m)
			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.verify != nil {
				tt.verify(t, m)
			}
		})
	}
}
