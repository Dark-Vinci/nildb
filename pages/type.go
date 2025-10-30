package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/interfaces"
)

type PageTypeConversion interface {
	faces.PageHandle
	FromPageHeader(header *bufferwheader.BufferWithHeader[PageHeader]) faces.PageHandle
	FromOverflowPageHeader(header *bufferwheader.BufferWithHeader[OverflowPageHeader]) faces.PageHandle
	FromDBHeader(header *bufferwheader.BufferWithHeader[DBHeader]) faces.PageHandle
}
