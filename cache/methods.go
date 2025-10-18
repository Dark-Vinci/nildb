package cache

import "github.com/dark-vinci/nildb/base"

func (c *Cache) GetMaxSize() uint {
	return c.MaxSize
}

func (c *Cache) GetPageSize() uint {
	return c.PageSize
}

func (c *Cache) Contains(pageNumber base.PageNumber) bool {
	_, exists := c.Pages[pageNumber]
	return exists
}

func (c *Cache) refPage(pageNumber base.PageNumber) *base.FrameID {
	if frameID, exists := c.Pages[pageNumber]; exists {
		c.updateHistory(frameID)
		return &frameID
	}

	return nil
}

func (c *Cache) Get(pageNumber base.PageNumber) *base.FrameID {
	return c.refPage(pageNumber)
}
