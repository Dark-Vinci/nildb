package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/interfaces"
)

type PageHeader struct {
	id uint32
}

// Page B+TREE PAGE
type Page struct {
	buffer *bufferwheader.BufferWithHeader[PageHeader]
}

var _ PageTypeConversion = (*Page)(nil)

func (p *Page) IsOverflow() bool {
	return false
}

func (p *Page) IntoBuffer() interface{} {
	return p.buffer
}

func (p *Page) FromPageHeader(header *bufferwheader.BufferWithHeader[PageHeader]) interfaces.PageHandle {
	return &Page{buffer: header}
}

func (p *Page) FromOverflowPageHeader(buffer *bufferwheader.BufferWithHeader[OverflowPageHeader]) interfaces.PageHandle {
	return &Page{
		buffer: bufferwheader.NewBufferWithHeader[PageHeader](buffer.Size()),
	}
}

func (p *Page) FromDBHeader(buffer *bufferwheader.BufferWithHeader[DBHeader]) interfaces.PageHandle {
	return &Page{
		buffer: bufferwheader.NewBufferWithHeader[PageHeader](buffer.Size()),
	}
}
