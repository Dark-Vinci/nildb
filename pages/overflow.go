package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/interfaces"
)

type OverflowPageHeader struct {
	offset uint64
}

type OverflowPage struct {
	buffer *bufferwheader.BufferWithHeader[OverflowPageHeader]
}

var _ PageTypeConversion = (*OverflowPage)(nil)

func (o *OverflowPage) IsOverflow() bool {
	return true
}

func (o *OverflowPage) FromBuffer(buffer []byte) interfaces.PageHandle {
	o.buffer = bufferwheader.FromSlice[OverflowPageHeader](buffer)

	return o
}

func (o *OverflowPage) IntoBuffer() (interface{}, error) {
	return o.buffer, nil
}

func (o *OverflowPage) FromPageHeader(buffer *bufferwheader.BufferWithHeader[PageHeader]) interfaces.PageHandle {
	return &OverflowPage{
		buffer: bufferwheader.NewBufferWithHeader[OverflowPageHeader](buffer.Size()),
	}
}

func (o *OverflowPage) FromOverflowPageHeader(buffer *bufferwheader.BufferWithHeader[OverflowPageHeader]) interfaces.PageHandle {
	return &OverflowPage{buffer: buffer}
}

func (o *OverflowPage) FromDBHeader(buffer *bufferwheader.BufferWithHeader[DBHeader]) interfaces.PageHandle {
	return &OverflowPage{
		buffer: bufferwheader.NewBufferWithHeader[OverflowPageHeader](buffer.Size()),
	}
}

func (o *OverflowPage) Type() string {
	return constants.OverFlowPage
}
