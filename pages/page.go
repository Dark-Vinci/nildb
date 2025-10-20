package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/constants"
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

func (p *Page) FromBuffer(buffer []byte) interfaces.PageHandle {
	p.buffer = bufferwheader.FromSlice[PageHeader](buffer)

	return p
}

func (p *Page) IntoBuffer() (interface{}, error) {
	return p.buffer, nil
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

func (p *Page) Type() string {
	return constants.BTreePage
}
