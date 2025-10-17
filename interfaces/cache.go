package interfaces

type Cache interface {
	Evict() (int, bool)
	Remove(frameID int)
	Size() int
	SetEvictable(frameID int, evictable bool)
	RecordAccess(frameID int)
}
