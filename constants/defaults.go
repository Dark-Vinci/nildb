package constants

const (
	DefaultMaxCacheSize       = 1000      // Example value
	MinCacheSize              = 10        // Example value
	DefaultPageSize           = 4096      // Example value
	DefaultPinPercentageLimit = 50.0      // Example value
	DefaultLruK               = 2         // Default K for LRU-K
	DefaultCRP                = uint64(0) // Default correlated reference period (logical accesses)
	DirtyFlag                 = 0x02
	PinnedFlag                = 0x04
)
