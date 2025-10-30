package frame

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
)

type Frame struct {
	PageNumber base.PageNumber
	Page       faces.PageHandle
	History    []uint64
	Last       uint64
	Flags      uint8
}

func NewFrame(pageNumber base.PageNumber, page faces.PageHandle) *Frame {
	return &Frame{
		PageNumber: pageNumber,
		Page:       page,
		History:    nil,
		Last:       0,
		Flags:      0,
	}
}

// Set sets the specified flags
func (f *Frame) Set(flags uint8) {
	f.Flags |= flags
}

// Unset clears the specified flags
func (f *Frame) Unset(flags uint8) {
	f.Flags &^= flags
}

// IsSet checks if the specified flags are set
func (f *Frame) IsSet(flags uint8) bool {
	return f.Flags&flags != 0
}

// IsOverflow checks if the page is an overflow page
func (f *Frame) IsOverflow() bool {
	return f.Page.IsOverflow()
}
