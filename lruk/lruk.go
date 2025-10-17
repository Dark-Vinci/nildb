package lruk

import (
	"sync"

	"github.com/dark-vinci/nildb/interfaces"
)

type LRUKCache struct {
	locker        sync.Mutex
	k             int
	capacity      int
	accessHistory map[int][]int64
	evictable     map[int]bool
}

func NewLRUKCache(capacity int, k int) *LRUKCache {
	return &LRUKCache{
		k:             k,
		capacity:      capacity,
		accessHistory: make(map[int][]int64),
		evictable:     make(map[int]bool),
	}
}

var _ interfaces.Cache = (*LRUKCache)(nil)
