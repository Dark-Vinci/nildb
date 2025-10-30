package pager

import (
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/interfaces"
)

type Builder struct {
	BlockSize uint32
	PageSize  uint32
	Cache     *faces.Cache
}

func NewBuilder() *Builder {
	return &Builder{
		BlockSize: 0,
		PageSize:  uint32(constants.DefaultPageSize),
		Cache:     nil,
	}
}

func (b *Builder) SetBlockSize(blockSize uint32) *Builder {
	b.BlockSize = blockSize
	return b
}

func (b *Builder) SetPageSize(pageSize uint32) *Builder {
	b.PageSize = pageSize
	return b
}

func (b *Builder) SetCache(cache faces.Cache) *Builder {
	b.Cache = &cache
	return b
}

// WE NEED TO SET CACHE PAGE SIZE TO Pager Page size
