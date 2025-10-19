package constants

import "time"

const (
	DefaultMaxCacheSize       = 1000
	MinCacheSize              = 10
	DefaultPageSize           = 4096
	DefaultPinPercentageLimit = 50.0
	DefaultLruK               = 2
	DefaultCRP                = uint64(0)
	DirtyFlag                 = 0x02
	PinnedFlag                = 0x04
	BatchSize                 = 100
	BatchTimeout              = time.Millisecond * 10
	PageAlignment             = 4096
	MinPageSize               = 512   // Minimum page size.
	MaxPageSize               = 65536 // Maximum page size.
	CellAlignment             = 8
)
