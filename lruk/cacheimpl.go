package lruk

import (
	"time"
)

func (l *LRUKCache) RecordAccess(frameID int) {
	l.locker.Lock()
	defer l.locker.Unlock()

	now := time.Now().UnixNano()

	l.accessHistory[frameID] = append(l.accessHistory[frameID], now)

	if len(l.accessHistory[frameID]) > l.k {
		l.accessHistory[frameID] = l.accessHistory[frameID][len(l.accessHistory[frameID])-l.k:]
	}
}

func (l *LRUKCache) Remove(frameID int) {
	l.locker.Lock()
	defer l.locker.Unlock()

	delete(l.accessHistory, frameID)
	delete(l.evictable, frameID)
}

func (l *LRUKCache) Size() int {
	l.locker.Lock()
	defer l.locker.Unlock()

	count := 0

	for _, isEvictable := range l.evictable {
		if isEvictable {
			count++
		}
	}

	return count
}
