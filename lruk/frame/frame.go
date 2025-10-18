package frame

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/interfaces"
)

type Frame struct {
	PageNumber base.PageNumber
	Page       *interfaces.RepPage
}

func (f *Frame) GetPage() *interfaces.RepPage {
	return f.Page
}

func NewFrame(pageNumber base.PageNumber, page *interfaces.RepPage) *Frame {
	return &Frame{
		PageNumber: pageNumber,
		Page:       page,
	}
}
