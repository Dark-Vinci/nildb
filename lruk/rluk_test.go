package lruk

import (
	"testing"
	"time"
)

func TestLRUK_BasicEviction(t *testing.T) {
	cache := NewLRUKCache(3, 2)

	cache.RecordAccess(1)
	time.Sleep(100 * time.Millisecond)

	cache.RecordAccess(2)

	cache.SetEvictable(1, true)
	cache.SetEvictable(2, true)

	victim, ok := cache.Evict()
	if !ok {
		t.Fatal("Expected a frame to be evicted")
	}

	if victim != 1 && victim != 2 {
		t.Fatalf("Unexpected victim: %d", victim)
	}
}
