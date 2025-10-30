package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/interfaces"
)

type DBHeader struct {
	version uint16
}

type PageZero struct {
	buffer *bufferwheader.BufferWithHeader[DBHeader]
}

var _ PageTypeConversion = (*PageZero)(nil)

func (p *PageZero) IsOverflow() bool {
	return p.asBtreePage().IsOverflow()
}

func (p *PageZero) asBtreePage() *Page {
	return &Page{
		buffer: bufferwheader.NewBufferWithHeader[PageHeader](p.buffer.Size()),
	}
}

func (p *PageZero) IntoBuffer() (interface{}, error) {
	return p.buffer, nil
}

func (p *PageZero) FromBuffer(buffer []byte) faces.PageHandle {
	p.buffer = bufferwheader.FromSlice[DBHeader](buffer)

	return p
}

func (p *PageZero) FromPageHeader(buffer *bufferwheader.BufferWithHeader[PageHeader]) faces.PageHandle {
	return &PageZero{
		buffer: bufferwheader.NewBufferWithHeader[DBHeader](buffer.Size()),
	}
}

func (p *PageZero) FromOverflowPageHeader(buffer *bufferwheader.BufferWithHeader[OverflowPageHeader]) faces.PageHandle {
	return &PageZero{
		buffer: bufferwheader.NewBufferWithHeader[DBHeader](buffer.Size()),
	}
}

func (p *PageZero) FromDBHeader(buffer *bufferwheader.BufferWithHeader[DBHeader]) faces.PageHandle {
	return &PageZero{buffer: buffer}
}

func (p *PageZero) Type() string {
	return constants.PageZero
}
