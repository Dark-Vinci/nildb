package cache

import (
	"math"

	"github.com/dark-vinci/nildb/base"
	"github.com/dark-vinci/nildb/constants"
	"github.com/dark-vinci/nildb/faces"
	"github.com/dark-vinci/nildb/frame"
	"github.com/dark-vinci/nildb/pages"
)

func (c *Cache) Load(pageNumber base.PageNumber, page *faces.PageHandle) *faces.PageHandle {
	if page == nil {
		return nil
	}

	frameID := c.Map(pageNumber)
	oldPage := c.Buffer[frameID].Page

	c.Buffer[frameID].Page = *page

	return &oldPage
}

// MustEvictDirtyPage checks if the next page to be evicted is dirty
func (c *Cache) MustEvictDirtyPage() bool {
	if len(c.Buffer) < int(c.MaxSize) {
		return false
	}

	victim := c.findVictim()

	return c.Buffer[victim].IsSet(constants.DirtyFlag)
}

func (c *Cache) findVictim() base.FrameID {
	var (
		t             = c.CurrentTime
		minVal        = uint64(math.MaxUint64)
		victim        = ^base.FrameID(0) // base.FrameID(uint64(math.MaxUint64))
		foundEligible = false
	)

	for id := 0; id < len(c.Buffer); id++ {
		var (
			fid = base.FrameID(id)
			fr  = c.Buffer[id]
		)

		if t-fr.Last <= c.CRP {
			continue
		}

		if !c.isEvictable(fid) {
			continue
		}

		historyK := fr.History[c.K-1]

		if historyK < minVal {
			minVal = historyK
			victim = fid
			foundEligible = true
		}
	}

	if foundEligible {
		return victim
	}

	minVal = uint64(math.MaxUint64)
	for id := 0; id < len(c.Buffer); id++ {
		fid := base.FrameID(id)
		fr := c.Buffer[id]

		if !c.isEvictable(fid) {
			continue
		}

		historyK := fr.History[c.K-1]
		if historyK < minVal {
			minVal = historyK
			victim = fid
		}
	}

	if victim == ^base.FrameID(0) {
		panic("No evictable frame found")
	}

	return victim
}

// isEvictable checks if a frame can be safely evicted
func (c *Cache) isEvictable(frameID base.FrameID) bool {
	fram := c.Buffer[frameID]

	return !fram.IsSet(constants.PinnedFlag) && !fram.IsOverflow()
}

// setFlags sets the specified flags for a page
func (c *Cache) setFlags(pageNumber base.PageNumber, flags uint8) bool {
	if frameId, exists := c.Pages[pageNumber]; exists {
		c.Buffer[frameId].Set(flags)
		return true
	}

	return false
}

func (c *Cache) Map(pageNumber base.PageNumber) base.FrameID {
	if frameID, exists := c.Pages[pageNumber]; exists {
		return frameID
	}

	var (
		frameID base.FrameID
		f       = &frame.Frame{}
	)

	// Buffer is not full, allocate a new page
	if uint(len(c.Buffer)) < c.MaxSize {
		frameID = base.FrameID(len(c.Buffer))
		f = frame.NewFrame(pageNumber, pages.Alloc(int(c.PageSize)))

		c.Buffer = append(c.Buffer, f)
	} else {
		// Buffer full, find a victim to evict
		victimID := c.findVictim()

		f = c.Buffer[victimID]

		delete(c.Pages, f.PageNumber)

		f.PageNumber = pageNumber
		f.Flags = 0
	}

	// Update history for the new or evicted frame
	c.updateHistory(frameID)

	c.Pages[pageNumber] = frameID
	return frameID
}

func (c *Cache) updateHistory(frameID base.FrameID) {
	fram := c.Buffer[frameID]

	c.CurrentTime++

	t := c.CurrentTime

	if len(fram.History) == 0 {
		fram.History = make([]uint64, c.K)
		fram.History[0] = t

		// history[1:] remains 0
		fram.Last = t

		return
	}

	if t-fram.Last > c.CRP {
		corrPeriod := fram.Last - fram.History[0]

		for i := int(c.K) - 1; i >= 1; i-- {
			fram.History[i] = fram.History[i-1] + corrPeriod
		}

		fram.History[0] = t
	}

	fram.Last = t
}

// unsetFlags clears the specified flags for a page
func (c *Cache) unsetFlags(pageNumber base.PageNumber, flags uint8) bool {
	if frameId, exists := c.Pages[pageNumber]; exists {
		c.Buffer[frameId].Unset(flags)
		return true
	}

	return false
}

// MarkDirty marks a page as dirty
func (c *Cache) MarkDirty(pageNumber base.PageNumber) bool {
	return c.setFlags(pageNumber, constants.DirtyFlag)
}

// MarkClean marks a page as clean
func (c *Cache) MarkClean(pageNumber base.PageNumber) bool {
	return c.unsetFlags(pageNumber, constants.DirtyFlag)
}

// Pin marks a page as unevictable
func (c *Cache) Pin(pageNumber base.PageNumber) bool {
	pinnedPercentage := float32(c.PinnedPages) / float32(c.MaxSize) * 100.0

	if pinnedPercentage >= c.PinPercentageLimit {
		return false
	}

	pinned := c.setFlags(pageNumber, constants.PinnedFlag)
	if pinned {
		c.PinnedPages++
	}

	return pinned
}

// Unpin marks a page as evictable
func (c *Cache) Unpin(pageNumber base.PageNumber) bool {
	unpinned := c.unsetFlags(pageNumber, constants.PinnedFlag)

	if unpinned {
		c.PinnedPages--
	}

	return unpinned
}

// Invalidate removes a page from the cache
func (c *Cache) Invalidate(pageNumber base.PageNumber) {
	if frameId, ok := c.Pages[pageNumber]; ok {
		c.Buffer[frameId].Flags = 0
		delete(c.Pages, pageNumber)
	}
}

// GetFrame retrieves the MemPage for a given FrameId
func (c *Cache) GetFrame(frameID base.FrameID) *faces.PageHandle {
	return &c.Buffer[frameID].Page
}
