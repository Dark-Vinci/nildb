package interfaces

import "github.com/dark-vinci/nildb/base"

type Cache interface {
	Evict() (int, bool)
	Remove(frameID int)
	Size() int
	SetEvictable(frameID int, evictable bool)
	RecordAccess(frameID int)

	Pin(pageNumber base.PageNumber) bool
	Unpin(pageNumber base.PageNumber) bool
	GetFrame(frameID base.FrameID) *RepPage
	MarkClean(pageNumber base.PageNumber) bool
	MarkDirty(pageNumber base.PageNumber) bool
	Map(pageNumber base.PageNumber) base.FrameID
	Load(pageNumber base.PageNumber, page *RepPage) *RepPage
	MustEvictDirtyPage() bool
	GetPage(pageNumber base.PageNumber) *uint

	RLock()
	RUnlock()
	Lock()
	Unlock()
	MaxSize() int
}
