package pages

import (
	"github.com/dark-vinci/nildb/bufferwheader"
	"github.com/dark-vinci/nildb/interfaces"
)

func Alloc(size int) faces.PageHandle {
	return &Page{buffer: bufferwheader.ForPage[PageHeader](size)}
}

func ReinitAs[T PageTypeConversion](memPage *faces.PageHandle) {
	var (
		current = *memPage
		newPage faces.PageHandle
		t       T
	)

	switch v := current.(type) {
	case *PageZero:
		newPage = t.FromDBHeader(v.buffer)
	case *OverflowPage:
		newPage = t.FromOverflowPageHeader(v.buffer)
	case *Page:
		newPage = t.FromPageHeader(v.buffer)
	default:
		panic("unknown PageHandle type")
	}

	*memPage = newPage
}
