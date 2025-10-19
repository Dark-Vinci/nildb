package cache

import (
	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/frame"
)

type Cache struct {
	Buffer             []*frame.Frame
	Pages              map[base.PageNumber]base.FrameID
	MaxSize            uint
	PageSize           uint
	PinPercentageLimit float32
	PinnedPages        uint
	K                  uint
	CRP                uint64
	CurrentTime        uint64
}

func NewCache() *Cache {
	return NewBuilder().Build()
}

func WithMaxSize(maxSize uint) *Cache {
	return NewBuilder().SetMaxSize(maxSize).Build()
}

func WithPageSize(pageSize uint) *Cache {
	return NewBuilder().SetPageSize(pageSize).Build()
}
