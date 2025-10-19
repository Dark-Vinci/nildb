package cache

import (
	"fmt"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/frame"
)

type Builder struct {
	MaxSize            uint
	PageSize           uint
	PinPercentageLimit float32
	lruK               uint
	crp                uint64
}

func NewBuilder() *Builder {
	return &Builder{
		MaxSize:            constants.DefaultMaxCacheSize,
		PageSize:           constants.DefaultPageSize,
		PinPercentageLimit: constants.DefaultPinPercentageLimit,
		lruK:               constants.DefaultLruK,
		crp:                constants.DefaultCRP,
	}
}

func (b *Builder) SetMaxSize(maxSize uint) *Builder {
	if maxSize < constants.MinCacheSize {
		panic(fmt.Sprintf("Buffer pool size must be at least %d", constants.MinCacheSize))
	}

	b.MaxSize = maxSize
	return b
}

func (b *Builder) SetPinPercentageLimit(pinPercentageLimit float32) *Builder {
	if pinPercentageLimit < 0.0 || pinPercentageLimit > 100.0 {
		panic(fmt.Sprintf("pin_percentage_limit must be a percentage (0..=100), got %f", pinPercentageLimit))
	}

	b.PinPercentageLimit = pinPercentageLimit
	return b
}

func (b *Builder) SetPageSize(pageSize uint) *Builder {
	b.PageSize = pageSize

	return b
}

func (b *Builder) LruK(k uint) *Builder {
	if k == 0 {
		panic("K must be at least 1")
	}

	b.lruK = k

	return b
}

func (b *Builder) CorrelatedReferencePeriod(crp uint64) *Builder {
	b.crp = crp
	return b
}

func (b *Builder) Build() *Cache {
	return &Cache{
		Buffer:             make([]*frame.Frame, 0, b.MaxSize),
		Pages:              make(map[base.PageNumber]base.FrameID, b.MaxSize),
		MaxSize:            b.MaxSize,
		PageSize:           b.PageSize,
		PinPercentageLimit: b.PinPercentageLimit,
		PinnedPages:        0,
		K:                  b.lruK,
		CRP:                b.crp,
		CurrentTime:        0,
	}
}
