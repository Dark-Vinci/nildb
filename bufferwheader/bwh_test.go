package bufferwheader

import (
	"testing"
	"unsafe"

	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/utils"
)

// Header types for testing.
type PageHeader struct {
	id uint32
}

type OverflowPageHeader struct {
	offset uint64
}

type DBHeader struct {
	version uint16
}

func Test_Alloc(t *testing.T) {
	testCases := []struct {
		name        string
		shouldPanic bool
		header      interface{}
		size        int
	}{
		{
			name:        "ValidSize_PageHeader",
			header:      PageHeader{},
			size:        4096,
			shouldPanic: false,
		},
		{
			name:        "ValidSize_DbHeader",
			header:      DBHeader{},
			size:        8192,
			shouldPanic: false,
		},
		{
			name:        "InsufficientSize_PageHeader",
			header:      PageHeader{},
			size:        utils.GetSize(PageHeader{}),
			shouldPanic: true,
		},
		{
			name:        "UnalignedSize",
			header:      PageHeader{},
			size:        4097,
			shouldPanic: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic")
					}
				}()
			}

			switch h := tt.header.(type) {
			case PageHeader:
				buffer := Allocate[PageHeader](tt.size)

				if !tt.shouldPanic {
					if len(buffer) != tt.size {
						t.Errorf("expected buffer size %d, got %d", tt.size, len(buffer))
					}

					//must all be zeroed buffer
					for i, b := range buffer {
						if b != 0 {
							t.Errorf("expected zeroed buffer, got non-zero at index %d: %v", i, b)
						}
					}
				}

			case DBHeader:
				buffer := Allocate[DBHeader](tt.size)
				if !tt.shouldPanic {
					if len(buffer) != tt.size {
						t.Errorf("expected buffer size %d, got %d", tt.size, len(buffer))
					}

					//must all be zeroed buffer
					for i, b := range buffer {
						if b != 0 {
							t.Errorf("expected zeroed buffer, got non-zero at index %d: %v", i, b)
						}
					}
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func Test_NewBufferWithHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      interface{}
		size        int
		shouldPanic bool
	}{
		{
			name:        "ValidSize_PageHeader",
			header:      PageHeader{},
			size:        4096,
			shouldPanic: false,
		},
		{
			name:        "ValidSize_DbHeader",
			header:      DBHeader{},
			size:        8192,
			shouldPanic: false,
		},
		{
			name:        "InsufficientSize_PageHeader",
			header:      PageHeader{},
			size:        utils.GetSize(PageHeader{}),
			shouldPanic: true,
		},
		{
			name:        "UnalignedSize",
			header:      PageHeader{},
			size:        4097,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic")
					}
				}()
			}

			switch h := tt.header.(type) {
			case PageHeader:
				buf := NewBufferWithHeader[PageHeader](tt.size)

				if !tt.shouldPanic {
					if buf.size != tt.size {
						t.Errorf("expected size %d, got %d", tt.size, buf.size)
					}

					headerSize := utils.GetSize(PageHeader{})
					if len(buf.Content()) != tt.size-headerSize {
						t.Errorf("expected content length %d, got %d", tt.size-headerSize, len(buf.Content()))
					}

					if buf.Header() == nil {
						t.Error("expected non-nil header")
					}

					if buf.Header().id != 0 {
						t.Errorf("expected header id 0, got %d", buf.Header().id)
					}
				}

			case DBHeader:
				buf := NewBufferWithHeader[DBHeader](tt.size)

				if !tt.shouldPanic {
					if buf.size != tt.size {
						t.Errorf("expected size %d, got %d", tt.size, buf.size)
					}

					headerSize := utils.GetSize(DBHeader{})
					if len(buf.Content()) != tt.size-headerSize {
						t.Errorf("expected content length %d, got %d", tt.size-headerSize, len(buf.Content()))
					}

					if buf.Header() == nil {
						t.Error("expected non-nil header")
					}

					if buf.Header().version != 0 {
						t.Errorf("expected header version 0, got %d", buf.Header().version)
					}
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestForPage(t *testing.T) {
	tests := []struct {
		name        string
		header      interface{}
		size        int
		shouldPanic bool
	}{
		{
			name:        "ValidSize",
			header:      PageHeader{},
			size:        4096,
			shouldPanic: false,
		},
		{
			name:        "MinSize",
			header:      DBHeader{},
			size:        constants.MinPageSize,
			shouldPanic: false,
		},
		{
			name:        "MaxSize",
			header:      PageHeader{},
			size:        constants.MaxPageSize,
			shouldPanic: false,
		},
		{
			name:        "BelowMinSize",
			header:      PageHeader{},
			size:        constants.MaxPageSize - 1,
			shouldPanic: true,
		},
		{
			name:        "AboveMaxSize",
			header:      DBHeader{},
			size:        constants.MinPageSize + 1,
			shouldPanic: true,
		},
		{
			name:        "UnalignedSize",
			header:      PageHeader{},
			size:        4097,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic")
					}
				}()
			}

			switch h := tt.header.(type) {
			case PageHeader:
				buf := ForPage[PageHeader](tt.size)
				if !tt.shouldPanic {
					if buf.size != tt.size {
						t.Errorf("expected size %d, got %d", tt.size, buf.size)
					}
				}

			case DBHeader:
				buf := ForPage[DBHeader](tt.size)

				if !tt.shouldPanic {
					if buf.size != tt.size {
						t.Errorf("expected size %d, got %d", tt.size, buf.size)
					}
				}
			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name        string
		header      interface{}
		buffer      []byte
		shouldPanic bool
	}{
		{
			name:        "ValidBuffer_PageHeader",
			header:      PageHeader{},
			buffer:      make([]byte, 4096),
			shouldPanic: false,
		},
		{
			name:        "ValidBuffer_DbHeader",
			header:      DBHeader{},
			buffer:      make([]byte, 4096),
			shouldPanic: false,
		},
		{
			name:        "InsufficientBufferSize",
			header:      PageHeader{},
			buffer:      make([]byte, utils.GetSize(PageHeader{})),
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic")
					}
				}()
			}

			switch h := tt.header.(type) {
			case PageHeader:
				buf := FromSlice[PageHeader](tt.buffer)
				if !tt.shouldPanic {
					if buf.size != len(tt.buffer) {
						t.Errorf("expected size %d, got %d", len(tt.buffer), buf.size)
					}

					headerSize := utils.GetSize(PageHeader{})
					if len(buf.Content()) != len(tt.buffer)-headerSize {
						t.Errorf("expected content length %d, got %d", len(tt.buffer)-headerSize, len(buf.Content()))
					}

					if buf.Header() == nil {
						t.Error("expected non-nil header")
					}
				}

			case DBHeader:
				buf := FromSlice[DBHeader](tt.buffer)
				if !tt.shouldPanic {
					if buf.size != len(tt.buffer) {
						t.Errorf("expected size %d, got %d", len(tt.buffer), buf.size)
					}

					headerSize := utils.GetSize(DBHeader{})
					if len(buf.Content()) != len(tt.buffer)-headerSize {
						t.Errorf("expected content length %d, got %d", len(tt.buffer)-headerSize, len(buf.Content()))
					}

					if buf.Header() == nil {
						t.Error("expected non-nil header")
					}
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestCast(t *testing.T) {
	tests := []struct {
		name        string
		fromHeader  interface{}
		toHeader    interface{}
		size        int
		shouldPanic bool
	}{
		{
			name:        "PageHeaderToDbHeader",
			fromHeader:  PageHeader{},
			toHeader:    DBHeader{},
			size:        4096,
			shouldPanic: false,
		},
		{
			name:        "DbHeaderToOverflowPageHeader",
			fromHeader:  DBHeader{},
			toHeader:    OverflowPageHeader{},
			size:        4096,
			shouldPanic: false,
		},
		{
			name:        "InsufficientSize",
			fromHeader:  PageHeader{},
			toHeader:    OverflowPageHeader{},
			size:        utils.GetSize(OverflowPageHeader{}),
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic")
					}
				}()
			}

			switch fromH := tt.fromHeader.(type) {
			case PageHeader:
				orig := NewBufferWithHeader[PageHeader](tt.size)

				switch toH := tt.toHeader.(type) {
				case DBHeader:
					newBuf := Cast[PageHeader, DBHeader](orig)

					if !tt.shouldPanic {
						if newBuf.size != orig.size {
							t.Errorf("expected size %d, got %d", orig.size, newBuf.size)
						}

						headerSize := utils.GetSize(DBHeader{})
						if len(newBuf.Content()) != orig.size-headerSize {
							t.Errorf("expected content length %d, got %d", orig.size-headerSize, len(newBuf.Content()))
						}

						if newBuf.Header() == nil {
							t.Error("expected non-nil header")
						}
					}

				case OverflowPageHeader:
					newBuf := Cast[PageHeader, OverflowPageHeader](orig)
					if !tt.shouldPanic {
						if newBuf.size != orig.size {
							t.Errorf("expected size %d, got %d", orig.size, newBuf.size)
						}

						headerSize := utils.GetSize(OverflowPageHeader{})
						if len(newBuf.Content()) != orig.size-headerSize {
							t.Errorf("expected content length %d, got %d", orig.size-headerSize, len(newBuf.Content()))
						}

						if newBuf.Header() == nil {
							t.Error("expected non-nil header")
						}
					}

				default:
					t.Fatalf("unsupported to header type: %T", toH)
				}
			case DBHeader:
				orig := NewBufferWithHeader[DBHeader](tt.size)

				switch toH := tt.toHeader.(type) {
				case OverflowPageHeader:
					newBuf := Cast[DBHeader, OverflowPageHeader](orig)

					if !tt.shouldPanic {
						if newBuf.size != orig.size {
							t.Errorf("expected size %d, got %d", orig.size, newBuf.size)
						}

						headerSize := utils.GetSize(OverflowPageHeader{})

						if len(newBuf.Content()) != orig.size-headerSize {
							t.Errorf("expected content length %d, got %d", orig.size-headerSize, len(newBuf.Content()))
						}

						if newBuf.Header() == nil {
							t.Error("expected non-nil header")
						}
					}

				default:
					t.Fatalf("unsupported to header type: %T", toH)
				}
			default:
				t.Fatalf("unsupported from header type: %T", fromH)
			}
		})
	}
}

func TestUsableSpace(t *testing.T) {
	tests := []struct {
		name     string
		header   interface{}
		size     int
		expected uint16
	}{
		{
			name:     "PageHeader",
			header:   PageHeader{},
			size:     4096,
			expected: uint16(4096 - utils.GetSize(PageHeader{})),
		},
		{
			name:     "DbHeader",
			header:   DBHeader{},
			size:     4096,
			expected: uint16(4096 - utils.GetSize(DBHeader{})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch h := tt.header.(type) {
			case PageHeader:
				buf := NewBufferWithHeader[PageHeader](tt.size)
				if got := buf.UsableSpace(); got != tt.expected {
					t.Errorf("expected usable space %d, got %d", tt.expected, got)
				}

			case DBHeader:
				buf := NewBufferWithHeader[DBHeader](tt.size)
				if got := buf.UsableSpace(); got != tt.expected {
					t.Errorf("expected usable space %d, got %d", tt.expected, got)
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	tests := []struct {
		name   string
		header interface{}
		size   int
	}{
		{
			name:   "PageHeader",
			header: PageHeader{},
			size:   4096,
		},
		{
			name:   "DBHeader",
			header: DBHeader{},
			size:   4096,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch h := tt.header.(type) {
			case PageHeader:
				buf := NewBufferWithHeader[PageHeader](tt.size)
				if buf.Header() == nil {
					t.Error("expected non-nil header")
				}

				if buf.Header().id != 0 {
					t.Errorf("expected header id 0, got %d", buf.Header().id)
				}

			case DBHeader:
				buf := NewBufferWithHeader[DBHeader](tt.size)
				if buf.Header() == nil {
					t.Error("expected non-nil header")
				}

				if buf.Header().version != 0 {
					t.Errorf("expected header version 0, got %d", buf.Header().version)
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestContent(t *testing.T) {
	tests := []struct {
		name   string
		header interface{}
		size   int
	}{
		{
			name:   "PageHeader",
			header: PageHeader{},
			size:   4096,
		},
		{
			name:   "DBHeader",
			header: DBHeader{},
			size:   4096,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch h := tt.header.(type) {
			case PageHeader:
				buf := NewBufferWithHeader[PageHeader](tt.size)
				headerSize := utils.GetSize(PageHeader{})

				if len(buf.Content()) != tt.size-headerSize {
					t.Errorf("expected content length %d, got %d", tt.size-headerSize, len(buf.Content()))
				}

				for i, b := range buf.Content() {
					if b != 0 {
						t.Errorf("expected zeroed content, got non-zero at index %d: %v", i, b)
					}
				}

			case DBHeader:
				buf := NewBufferWithHeader[DBHeader](tt.size)
				headerSize := utils.GetSize(DBHeader{})

				if len(buf.Content()) != tt.size-headerSize {
					t.Errorf("expected content length %d, got %d", tt.size-headerSize, len(buf.Content()))
				}

				for i, b := range buf.Content() {
					if b != 0 {
						t.Errorf("expected zeroed content, got non-zero at index %d: %v", i, b)
					}
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}

func TestAsSlice(t *testing.T) {
	tests := []struct {
		name   string
		header interface{}
		size   int
	}{
		{
			name:   "PageHeader",
			header: PageHeader{},
			size:   4096,
		},
		{
			name:   "DBHeader",
			header: DBHeader{},
			size:   4096,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch h := tt.header.(type) {
			case PageHeader:
				buf := NewBufferWithHeader[PageHeader](tt.size)
				slice := buf.AsSlice()

				if len(slice) != tt.size {
					t.Errorf("expected slice length %d, got %d", tt.size, len(slice))
				}

				headerSize := utils.GetSize(PageHeader{})
				headerSlice := slice[:headerSize]
				headerPtr := (*PageHeader)(unsafe.Pointer(&headerSlice[0]))
				buf.Header().id = 42

				if headerPtr.id != 42 {
					t.Errorf("expected slice header id 42, got %d", headerPtr.id)
				}

			case DBHeader:
				buf := NewBufferWithHeader[DBHeader](tt.size)
				slice := buf.AsSlice()
				if len(slice) != tt.size {
					t.Errorf("expected slice length %d, got %d", tt.size, len(slice))
				}

				headerSize := utils.GetSize(DBHeader{})
				headerSlice := slice[:headerSize]
				headerPtr := (*DBHeader)(unsafe.Pointer(&headerSlice[0]))
				buf.Header().version = 123

				if headerPtr.version != 123 {
					t.Errorf("expected slice header version 123, got %d", headerPtr.version)
				}

			default:
				t.Fatalf("unsupported header type: %T", h)
			}
		})
	}
}
