package cache

import (
	"reflect"
	"testing"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/constants"
)

// TestMapAndGet verifies mapping and retrieving pages
func TestMapAndGet(t *testing.T) {
	c := NewBuilder().SetMaxSize(3).LruK(2).Build()

	// Map a page
	frameId := c.Map(1)
	if !c.Contains(1) {
		t.Errorf("Expected page 1 to be cached")
	}
	if got := c.Get(1); got == nil || *got != frameId {
		t.Errorf("Expected frame ID %d for page 1, got %v", frameId, got)
	}

	// Verify history update
	frame := c.Buffer[frameId]
	if len(frame.History) != 2 || frame.History[0] != 1 || frame.History[1] != 0 {
		t.Errorf("Expected history [1, 0], got %v", frame.History)
	}

	// Map another page and get mutable
	frameId2 := c.Map(2)
	if got := c.Get(2); got == nil || *got != frameId2 {
		t.Errorf("Expected frame ID %d for page 2, got %v", frameId2, got)
	}
	if !frame.IsSet(constants.DirtyFlag) {
		t.Errorf("Expected page 2 to be marked dirty")
	}
}

// TestLRUKEviction verifies LRU-K eviction policy
func TestLRUKEviction(t *testing.T) {
	c := NewBuilder().SetMaxSize(3).LruK(2).CorrelatedReferencePeriod(10).Build()

	// Fill cache
	c.Map(1) // Time 1
	c.Map(2) // Time 2
	c.Map(3) // Time 3

	// Access pages to set history
	c.Get(1) // Time 4
	c.Get(2) // Time 5
	c.Get(3) // Time 6

	// Verify history
	if c.Buffer[c.Pages[1]].History[0] != 4 || c.Buffer[c.Pages[1]].History[1] != 1 {
		t.Errorf("Expected history for page 1 [4, 1], got %v", c.Buffer[c.Pages[1]].History)
	}

	// Map a new page, should evict page 1 (oldest K-th access)
	c.CurrentTime = 20 // Ensure outside CRP
	frameId := c.Map(4)
	if c.Contains(1) {
		t.Errorf("Expected page 1 to be evicted")
	}
	if !c.Contains(4) {
		t.Errorf("Expected page 4 to be cached")
	}
	if c.Buffer[frameId].PageNumber != 4 {
		t.Errorf("Expected frame %d to hold page 4, got %d", frameId, c.Buffer[frameId].PageNumber)
	}
}

// TestCRPEviction verifies CRP handling in eviction
func TestCRPEviction(t *testing.T) {
	c := NewBuilder().SetMaxSize(3).LruK(2).CorrelatedReferencePeriod(5).Build()

	// Fill cache
	c.Map(1) // Time 1
	c.Map(2) // Time 2
	c.Map(3) // Time 3

	// Recent accesses within CRP
	c.Get(1) // Time 4
	c.Get(2) // Time 5

	// Try to map a new page, no eviction due to CRP
	c.CurrentTime = 7
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to no evictable pages within CRP")
		}
	}()
	c.Map(4)
}

// TestPinning verifies pinning and unpinning pages
func TestPinning(t *testing.T) {
	c := NewBuilder().
		SetMaxSize(2).
		SetPinPercentageLimit(50.0).
		Build()

	// Map and pin a page
	c.Map(1)
	if !c.Pin(1) {
		t.Errorf("Expected page 1 to be pinned")
	}
	if c.PinnedPages != 1 {
		t.Errorf("Expected 1 pinned page, got %d", c.PinnedPages)
	}
	if !c.Buffer[c.Pages[1]].IsSet(constants.PinnedFlag) {
		t.Errorf("Expected page 1 to have pinned flag set")
	}

	// Try to pin beyond limit
	c.Map(2)
	if c.Pin(2) {
		t.Errorf("Expected pinning page 2 to fail due to pin limit")
	}

	// Unpin page
	if !c.Unpin(1) {
		t.Errorf("Expected page 1 to be unpinned")
	}
	if c.PinnedPages != 0 {
		t.Errorf("Expected 0 pinned pages, got %d", c.PinnedPages)
	}
	if c.Buffer[c.Pages[1]].IsSet(constants.PinnedFlag) {
		t.Errorf("Expected page 1 to have pinned flag unset")
	}
}

// TestDirtyPages verifies dirty flag management
func TestDirtyPages(t *testing.T) {
	c := NewBuilder().SetPageSize(2).Build()

	// Map and mark dirty
	c.Map(1)
	if !c.MarkDirty(1) {
		t.Errorf("Expected page 1 to be marked dirty")
	}
	if !c.Buffer[c.Pages[1]].IsSet(constants.DirtyFlag) {
		t.Errorf("Expected page 1 to have dirty flag set")
	}

	// Mark clean
	if !c.MarkClean(1) {
		t.Errorf("Expected page 1 to be marked clean")
	}
	if c.Buffer[c.Pages[1]].IsSet(constants.DirtyFlag) {
		t.Errorf("Expected page 1 to have dirty flag unset")
	}

	// Check must evict dirty page
	c.MarkDirty(1)
	c.Map(2)
	if !c.MustEvictDirtyPage() {
		t.Errorf("Expected must evict dirty page to return true")
	}
}

// TestInvalidate verifies page invalidation
func TestInvalidate(t *testing.T) {
	c := NewBuilder().SetMaxSize(2).Build()

	// Map and invalidate a page
	c.Map(1)
	c.Pin(1)
	c.MarkDirty(1)
	c.Invalidate(1)
	if c.Contains(1) {
		t.Errorf("Expected page 1 to be invalidated")
	}
	if c.PinnedPages != 0 {
		t.Errorf("Expected 0 pinned pages after invalidation, got %d", c.PinnedPages)
	}
	if c.Buffer[c.Pages[1]].Flags != 0 {
		t.Errorf("Expected flags to be reset after invalidation, got %d", c.Buffer[c.Pages[1]].Flags)
	}
}

// TestGetManyMut verifies retrieving multiple mutable pages
func TestGetManyMut(t *testing.T) {
	c := NewBuilder().SetMaxSize(3).Build()

	// Map pages
	c.Map(1)
	c.Map(2)
	c.Map(3)

	// Get multiple mutable pages
	pages := c.GetMany([]PageNumber{1, 2})
	if len(pages) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(pages))
	}
	for _, pn := range []base.PageNumber{1, 2} {
		if !c.Buffer[c.Pages[pn]].IsSet(constants.DirtyFlag) {
			t.Errorf("Expected page %d to be marked dirty", pn)
		}
	}

	// Try to get non-existent page
	pages = c.GetManyMut([]base.PageNumber{1, 4})
	if pages != nil {
		t.Errorf("Expected nil for non-existent page, got %v", pages)
	}
}

// TestLoad verifies loading a page into the cache
func TestLoad(t *testing.T) {
	c := NewBuilder().SetMaxSize(2).Build()

	// Map a page
	frameId := c.Map(1)
	oldPage := c.GetFrame(frameId)

	// Load a new page
	newPage := MemPage{Data: []byte{1, 2, 3}}
	evicted := c.Load(2, newPage)
	if !reflect.DeepEqual(evicted, oldPage) {
		t.Errorf("Expected evicted page %v, got %v", oldPage, evicted)
	}

	if !reflect.DeepEqual(c.GetFrame(c.Pages[2]).Data, newPage.Data) {
		t.Errorf("Expected loaded page %v, got %v", newPage.Data, c.GetFrame(c.Pages[2]).Data)
	}
}
