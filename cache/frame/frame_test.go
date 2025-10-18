package frame

import (
	"testing"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/interfaces"
)

// MockRepPage is a mock implementation of interfaces.RepPage for testing
type MockRepPage struct {
	overflow bool
}

func (m *MockRepPage) Type() string {
	return "MockRepPage"
}

func (m *MockRepPage) IsOverFlow() bool {
	return m.overflow
}

// TestNewFrame verifies the NewFrame constructor
func TestNewFrame(t *testing.T) {
	pageNumber := base.PageNumber(42)
	page := &MockRepPage{overflow: false}

	frame := NewFrame(pageNumber, page)

	if frame.PageNumber != pageNumber {
		t.Errorf("Expected PageNumber %d, got %d", pageNumber, frame.PageNumber)
	}
	if frame.Page != page {
		t.Errorf("Expected Page %v, got %v", page, frame.Page)
	}
	if frame.History != nil {
		t.Errorf("Expected History to be nil, got %v", frame.History)
	}
	if frame.Last != 0 {
		t.Errorf("Expected Last to be 0, got %d", frame.Last)
	}
	if frame.Flags != 0 {
		t.Errorf("Expected Flags to be 0, got %d", frame.Flags)
	}
}

// TestSet verifies the Set method
func TestSet(t *testing.T) {
	frame := NewFrame(base.PageNumber(1), &MockRepPage{overflow: false})

	// Test setting single flag
	frame.Set(constants.DirtyFlag)
	if frame.Flags != constants.DirtyFlag {
		t.Errorf("Expected Flags to be %d, got %d", constants.DirtyFlag, frame.Flags)
	}

	// Test setting multiple flags
	frame.Set(constants.PinnedFlag)
	expected := uint8(constants.DirtyFlag | constants.PinnedFlag)
	if frame.Flags != expected {
		t.Errorf("Expected Flags to be %d, got %d", expected, frame.Flags)
	}

	// Test setting already set flag
	frame.Set(constants.DirtyFlag)
	if frame.Flags != expected {
		t.Errorf("Expected Flags to remain %d, got %d", expected, frame.Flags)
	}
}

// TestUnset verifies the Unset method
func TestUnset(t *testing.T) {
	frame := NewFrame(base.PageNumber(1), &MockRepPage{overflow: false})
	frame.Flags = constants.DirtyFlag | constants.PinnedFlag

	// Test unsetting single flag
	frame.Unset(constants.DirtyFlag)
	if frame.Flags != constants.PinnedFlag {
		t.Errorf("Expected Flags to be %d, got %d", constants.PinnedFlag, frame.Flags)
	}

	// Test unsetting multiple flags
	frame.Unset(constants.PinnedFlag)
	if frame.Flags != 0 {
		t.Errorf("Expected Flags to be 0, got %d", frame.Flags)
	}

	// Test unsetting non-set flag
	frame.Unset(constants.DirtyFlag)
	if frame.Flags != 0 {
		t.Errorf("Expected Flags to remain 0, got %d", frame.Flags)
	}
}

// TestIsSet verifies the IsSet method
func TestIsSet(t *testing.T) {
	frame := NewFrame(base.PageNumber(1), &MockRepPage{overflow: false})
	frame.Flags = constants.DirtyFlag

	tests := []struct {
		name     string
		flags    uint8
		expected bool
	}{
		{
			name:     "Single flag set",
			flags:    constants.DirtyFlag,
			expected: true,
		},
		{
			name:     "Single flag not set",
			flags:    constants.PinnedFlag,
			expected: false,
		},
		{
			name:     "Multiple flags, one set",
			flags:    constants.DirtyFlag | constants.PinnedFlag,
			expected: true,
		},
		{
			name:     "Multiple flags, none set",
			flags:    constants.PinnedFlag | 0x08,
			expected: false,
		},
		{
			name:     "Zero flags",
			flags:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := frame.IsSet(tt.flags)

			if got != tt.expected {
				t.Errorf("Expected IsSet(%d) to return %v, got %v", tt.flags, tt.expected, got)
			}
		})
	}
}

// TestIsOverflow verifies the IsOverflow method
func TestIsOverflow(t *testing.T) {
	tests := []struct {
		name     string
		page     interfaces.RepPage
		expected bool
	}{
		{
			name:     "Overflow page",
			page:     &MockRepPage{overflow: true},
			expected: true,
		},
		{
			name:     "Non-overflow page",
			page:     &MockRepPage{overflow: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := NewFrame(base.PageNumber(1), tt.page)
			got := frame.IsOverflow()

			if got != tt.expected {
				t.Errorf("Expected IsOverflow to return %v, got %v", tt.expected, got)
			}
		})
	}
}
