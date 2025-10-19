package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/interfaces"
)

type PageTypeConversion interface {
	interfaces.PageHandle
	FromPageHeader(header *bufferwheader.BufferWithHeader[PageHeader]) interfaces.PageHandle
	FromOverflowPageHeader(header *bufferwheader.BufferWithHeader[OverflowPageHeader]) interfaces.PageHandle
	FromDBHeader(header *bufferwheader.BufferWithHeader[DBHeader]) interfaces.PageHandle
}
